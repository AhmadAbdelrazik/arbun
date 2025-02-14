package handlers

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/repository"
	"fmt"
	"net/http"
	"testing"
)

func productToPostProductInput(p repository.Product) postProductInput {
	return postProductInput{
		Name:        p.Name,
		Description: p.Description,
		Vendor:      p.Vendor,
		Amount:      p.AvailableAmount,
		Properties:  p.Properties,
	}
}

func TestPostProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	t.Run("valid insertion", func(t *testing.T) {
		tests := []struct {
			testName string
			Product  repository.Product
		}{
			{
				testName: "first product",
				Product: repository.Product{
					ID:              1,
					Name:            "product 1",
					Description:     "description of product 1",
					Vendor:          "vendor 1",
					AvailableAmount: 5,
					Version:         1,
					Properties: map[string]string{
						"size": "12",
					},
				},
			},
			{
				testName: "second product",
				Product: repository.Product{
					ID:              2,
					Name:            "product 2",
					Description:     "description of product 2",
					Vendor:          "vendor 1",
					AvailableAmount: 8,
					Version:         1,
					Properties: map[string]string{
						"size": "14",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.testName, func(t *testing.T) {
				res, err := ts.Post("/products", productToPostProductInput(tt.Product))
				assert.Nil(t, err)

				var responseBody struct {
					Product repository.Product `json:"product"`
				}

				err = ts.ReadResponseBody(res, &responseBody)
				assert.Nil(t, err)
				assert.Equal(t, responseBody.Product.String(), tt.Product.String())
			})
		}
	})

	t.Run("invalid insertion", func(t *testing.T) {
		t.Run("duplicate product", func(t *testing.T) {
			product := repository.Product{
				ID:              3,
				Name:            "product 2",
				Description:     "description of product 2",
				Vendor:          "vendor 1",
				AvailableAmount: 4,
				Version:         1,
				Properties: map[string]string{
					"size": "14",
				},
			}

			var responseBody struct {
				Error string `json:"error"`
			}

			res, err := ts.Post("/products", productToPostProductInput(product))
			assert.Nil(t, err)

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, res.StatusCode, http.StatusBadRequest)
			assert.Equal(t, responseBody.Error, "product already exists")
		})
		t.Run("invalid product", func(t *testing.T) {
			product := repository.Product{
				ID:              3,
				Name:            "product 2",
				Description:     "description of product 2",
				Vendor:          "vendor 2",
				AvailableAmount: 4,
				Version:         1,
				Properties: map[string]string{
					"size": "14",
				},
			}

			var responseBody struct {
				Error string `json:"error"`
			}

			res, err := ts.Post("/products", product)
			assert.Nil(t, err)

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, res.StatusCode, http.StatusBadRequest)
			assert.Equal(t, responseBody.Error, `body contains unknown key "ID"`)
		})
	})
}

func TestGetProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	product1 := repository.Product{
		ID:          1,
		Name:        "product 1",
		Description: "description of product 1",
		Vendor:      "vendor 1",
		Properties: map[string]string{
			"size": "11",
		},
		Version:         1,
		AvailableAmount: 4,
	}
	product2 := repository.Product{
		ID:          2,
		Name:        "product 2",
		Description: "description of product 2",
		Vendor:      "vendor 1",
		Properties: map[string]string{
			"size": "14",
		},
		Version:         1,
		AvailableAmount: 6,
	}

	ts.Post("/products", productToPostProductInput(product1))
	ts.Post("/products", productToPostProductInput(product2))

	t.Run("valid fetching", func(t *testing.T) {
		tests := []struct {
			testName string
			ID       int64
			product  repository.Product
		}{
			{
				testName: "product 1",
				ID:       1,
				product:  product1,
			},
			{
				testName: "product 2",
				ID:       2,
				product:  product2,
			},
		}

		for _, tt := range tests {
			var responseBody struct {
				Product repository.Product `json:"product"`
			}
			t.Run(tt.testName, func(t *testing.T) {
				res, err := ts.Get(fmt.Sprintf("/products/%v", tt.ID))
				assert.Nil(t, err)

				err = ts.ReadResponseBody(res, &responseBody)
				assert.Nil(t, err)
				assert.Equal(t, responseBody.Product.String(), tt.product.String())
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

func productToPatchProductInput(p repository.Product) patchProductInput {
	return patchProductInput{
		Name:            &p.Name,
		Description:     &p.Description,
		Vendor:          &p.Vendor,
		AvailableAmount: &p.AvailableAmount,
		Properties:      p.Properties,
	}
}

func TestPatchProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	product1 := repository.Product{
		ID:          1,
		Name:        "product 1",
		Description: "description of product 1",
		Vendor:      "vendor 1",
		Properties: map[string]string{
			"size": "11",
		},
		Version:         1,
		AvailableAmount: 4,
	}

	ts.Post("/products", productToPostProductInput(product1))

	t.Run("valid update", func(t *testing.T) {
		newProduct := product1
		newProduct.Name = "product 1 updated"
		newProduct.Description = "product 1 description updated"
		newProduct.AvailableAmount = 53
		newProduct.Version++

		reqBody := patchProductInput{
			Name:            &newProduct.Name,
			Description:     &newProduct.Description,
			AvailableAmount: &newProduct.AvailableAmount,
		}

		res, err := ts.Patch("/products/1", reqBody)
		assert.Nil(t, err)

		var responseBody struct {
			Product repository.Product `json:"product"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Product.String(), newProduct.String())
	})

	t.Run("invalid update", func(t *testing.T) {
		newProduct := product1
		newProduct.Name = "product 1 updated"
		newProduct.Description = "product 1 description updated"
		newProduct.AvailableAmount = 53
		newProduct.Version++

		reqBody := patchProductInput{
			Name:            &newProduct.Name,
			Description:     &newProduct.Description,
			AvailableAmount: &newProduct.AvailableAmount,
		}

		res, err := ts.Patch("/products/3", reqBody)
		assert.Nil(t, err)

		var responseBody struct {
			Error string `json:"error"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Error, "the requested resource could not be found")
	})
}

func TestDeleteProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	product1 := repository.Product{
		ID:          1,
		Name:        "product 1",
		Description: "description of product 1",
		Vendor:      "vendor 1",
		Properties: map[string]string{
			"size": "11",
		},
		Version:         1,
		AvailableAmount: 4,
	}

	ts.Post("/products", productToPostProductInput(product1))

	t.Run("valid delete", func(t *testing.T) {
		res, err := ts.Delete("/products/1")
		assert.Nil(t, err)

		var responseBody struct {
			Message string `json:"message"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Message, "deleted successfully")
	})

	t.Run("invalid deletion", func(t *testing.T) {
		res, err := ts.Delete("/products/3")
		assert.Nil(t, err)

		var responseBody struct {
			Error string `json:"error"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Error, "the requested resource could not be found")
	})
}
