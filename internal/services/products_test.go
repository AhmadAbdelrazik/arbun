package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"testing"
)

func TestProductInsertion(t *testing.T) {
	t.Run("valid insertion", func(t *testing.T) {
		model := models.NewModel()
		productService := newProductService(model)

		tests := []struct {
			TestName string
			product  domain.Product
		}{
			{
				TestName: "product 1",
				product: domain.Product{
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
				product: domain.Product{
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
				got, err := productService.InsertProduct(tt.product)
				assert.Nil(t, err)
				assert.Equal(t, got.String(), tt.product.String())
			})
		}
	})

	t.Run("invalid insertion", func(t *testing.T) {
		model := models.NewModel()
		productService := newProductService(model)
		firstProduct := domain.Product{
			Name:        "product 1",
			Description: "product 1 description",
			Vendor:      "vendor 1",
			Price:       23.99,
			Properties: map[string]string{
				"size": "12",
			},
			AvailableAmount: 4,
		}
		productService.InsertProduct(firstProduct)

		invalidProduct := domain.Product{
			Name:        "product 1",
			Description: "product 1 description",
			Vendor:      "vendor 1",
			Price:       23.99,
			Properties: map[string]string{
				"size": "12",
			},
			AvailableAmount: 6,
		}

		p, err := productService.InsertProduct(invalidProduct)

		assert.Err(t, err, models.ErrDuplicateProduct)
		assert.Equal(t, p.String(), domain.Product{}.String())

	})
}

func TestProductFetching(t *testing.T) {
	model := models.NewModel()
	productService := newProductService(model)

	product1 := domain.Product{
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

	product2 := domain.Product{
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

	productService.InsertProduct(product1)
	productService.InsertProduct(product2)

	t.Run("valid fetch", func(t *testing.T) {
		p, err := productService.GetProductByID(1)

		assert.Nil(t, err)
		assert.Equal(t, p.String(), product1.String())
	})

	t.Run("invalid fetch", func(t *testing.T) {
		p, err := productService.GetProductByID(5)

		assert.Err(t, err, models.ErrProductNotFound)
		assert.Equal(t, p.String(), domain.Product{}.String())
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
	model := models.NewModel()
	service := newProductService(model)

	product1 := domain.Product{
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

	service.InsertProduct(product1)
	newName := "product 1 new name"

	t.Run("valid update", func(t *testing.T) {
		product, err := service.UpdateProduct(UpdateProductParam{ID: 1, Name: &newName})

		assert.Nil(t, err)
		assert.Equal(t, product.Name, newName)
	})

	t.Run("invalid update", func(t *testing.T) {
		t.Run("not providing id", func(t *testing.T) {
			p, err := service.UpdateProduct(UpdateProductParam{Name: &newName})

			assert.Err(t, err, models.ErrProductNotFound)
			assert.Equal(t, p.String(), domain.Product{}.String())
		})
		t.Run("non existent product update", func(t *testing.T) {
			p, err := service.UpdateProduct(UpdateProductParam{ID: 4, Name: &newName})

			assert.Err(t, err, models.ErrProductNotFound)
			assert.Equal(t, p.String(), domain.Product{}.String())
		})
	})
}

func TestProductDeletion(t *testing.T) {
	model := models.NewModel()
	service := newProductService(model)
	product1 := domain.Product{
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

	service.InsertProduct(product1)

	t.Run("invalid deletion", func(t *testing.T) {
		err := service.DeleteProduct(4)
		assert.Err(t, err, models.ErrProductNotFound)
	})

	t.Run("valid deletion", func(t *testing.T) {
		err := service.DeleteProduct(1)
		assert.Nil(t, err)

		products, err := service.GetAllProducts()
		assert.Nil(t, err)
		assert.Equal(t, len(products), 0)
	})
}
