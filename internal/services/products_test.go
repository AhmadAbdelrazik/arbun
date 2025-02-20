package services

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/repository"
	"testing"
)

func TestProductInsertion(t *testing.T) {
	t.Run("valid insertion", func(t *testing.T) {
		model := repository.NewModel()
		productService := newProductService(model)

		tests := []struct {
			TestName string
			product  repository.Product
		}{
			{
				TestName: "product 1",
				product: repository.Product{
					Name:        "product 1",
					ID:          1,
					Description: "product 1 description",
					Vendor:      "vendor 1",
					Price:       23.99,
					Properties: map[string]string{
						"size": "12",
					},
					AvailableAmount: 4,
					Version:         1,
				},
			},
			{
				TestName: "product 2",
				product: repository.Product{
					Name:        "product 2",
					ID:          2,
					Description: "product 2 description",
					Vendor:      "vendor 1",
					Price:       23.99,
					Properties: map[string]string{
						"size": "12",
					},
					AvailableAmount: 6,
					Version:         1,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.TestName, func(t *testing.T) {
				params := productToInsertParam(tt.product)

				got, err := productService.InsertProduct(params)
				assert.Nil(t, err)
				assert.Equal(t, got.String(), tt.product.String())
			})
		}
	})

	t.Run("invalid insertion", func(t *testing.T) {
		model := repository.NewModel()
		productService := newProductService(model)
		firstProduct := repository.Product{
			Name:        "product 1",
			Description: "product 1 description",
			Vendor:      "vendor 1",
			Price:       23.99,
			Properties: map[string]string{
				"size": "12",
			},
			AvailableAmount: 4,
		}
		productService.InsertProduct(productToInsertParam(firstProduct))

		invalidProduct := repository.Product{
			Name:        "product 1",
			Description: "product 1 description",
			Vendor:      "vendor 1",
			Price:       23.99,
			Properties: map[string]string{
				"size": "12",
			},
			AvailableAmount: 6,
		}

		product, err := productService.InsertProduct(productToInsertParam(invalidProduct))

		assert.Err(t, err, repository.ErrDuplicateProduct)
		assert.Equal(t, product.String(), repository.Product{}.String())

	})
}

func TestProductFetching(t *testing.T) {
	model := repository.NewModel()
	productService := newProductService(model)

	product1 := repository.Product{
		Name:        "product 1",
		Description: "product 1 description",
		Vendor:      "vendor 1",
		Price:       23.99,
		Properties: map[string]string{
			"size": "12",
		},
		AvailableAmount: 4,
		Version:         1,
		ID:              1,
	}

	product2 := repository.Product{
		Name:        "product 2",
		Description: "product 2 description",
		Vendor:      "vendor 2",
		Price:       23.99,
		Properties: map[string]string{
			"size": "14",
		},
		AvailableAmount: 4,
		Version:         1,
		ID:              2,
	}

	productService.InsertProduct(productToInsertParam(product1))
	productService.InsertProduct(productToInsertParam(product2))

	t.Run("valid fetch", func(t *testing.T) {
		p, err := productService.GetProductByID(1)

		assert.Nil(t, err)
		assert.Equal(t, p.String(), product1.String())
	})

	t.Run("invalid fetch", func(t *testing.T) {
		p, err := productService.GetProductByID(5)

		assert.Err(t, err, repository.ErrProductNotFound)
		assert.Equal(t, p.String(), repository.Product{}.String())
	})

	t.Run("fetch All", func(t *testing.T) {
		products, err := productService.GetAllProducts()

		assert.Nil(t, err)
		assert.Equal(t, len(products), 2)
		assert.Equal(t, product1.String(), products[0].String())
		assert.Equal(t, product2.String(), products[1].String())
	})
}

func TestProductUpdate(t *testing.T) {
	model := repository.NewModel()
	service := newProductService(model)

	product1 := repository.Product{
		Name:        "product 1",
		Description: "product 1 description",
		Vendor:      "vendor 1",
		Price:       23.99,
		Properties: map[string]string{
			"size": "12",
		},
		AvailableAmount: 4,
		Version:         1,
		ID:              1,
	}

	service.InsertProduct(productToInsertParam(product1))
	newName := "product 1 new name"

	t.Run("valid update", func(t *testing.T) {
		product, err := service.UpdateProduct(UpdateProductParam{ID: 1, Name: &newName})

		assert.Nil(t, err)
		assert.Equal(t, product.Name, newName)
	})

	t.Run("invalid update", func(t *testing.T) {
		t.Run("not providing id", func(t *testing.T) {
			product, err := service.UpdateProduct(UpdateProductParam{Name: &newName})

			assert.Err(t, err, repository.ErrProductNotFound)
			assert.Equal(t, product.String(), repository.Product{}.String())
		})
		t.Run("non existent product update", func(t *testing.T) {
			product, err := service.UpdateProduct(UpdateProductParam{ID: 4, Name: &newName})

			assert.Err(t, err, repository.ErrProductNotFound)
			assert.Equal(t, product.String(), repository.Product{}.String())
		})
	})
}

func TestProductDeletion(t *testing.T) {
	model := repository.NewModel()
	service := newProductService(model)
	product1 := repository.Product{
		Name:        "product 1",
		Description: "product 1 description",
		Vendor:      "vendor 1",
		Price:       23.99,
		Properties: map[string]string{
			"size": "12",
		},
		AvailableAmount: 4,
		Version:         1,
		ID:              1,
	}

	service.InsertProduct(productToInsertParam(product1))

	t.Run("invalid deletion", func(t *testing.T) {
		err := service.DeleteProduct(4)
		assert.Err(t, err, repository.ErrProductNotFound)
	})

	t.Run("valid deletion", func(t *testing.T) {
		err := service.DeleteProduct(1)
		assert.Nil(t, err)

		products, err := service.GetAllProducts()
		assert.Nil(t, err)
		assert.Equal(t, len(products), 0)
	})
}

func productToInsertParam(p repository.Product) InsertProductParam {
	return InsertProductParam{
		Name:            p.Name,
		Description:     p.Description,
		Vendor:          p.Vendor,
		Properties:      p.Properties,
		AvailableAmount: p.AvailableAmount,
		Price:           p.Price,
	}
}
