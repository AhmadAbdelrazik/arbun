package services

import "AhmadAbdelrazik/arbun/internal/repository"

type ProductService struct {
	model repository.ProductModel
}

func (p ProductService) InsertProduct(name string, description string, properties map[string]string) (int64, error) {
	product := repository.Product{
		Name:        name,
		Description: description,
		Properties:  properties,
	}

	return p.model.InsertProduct(product)

}

func (p ProductService) GetProductByID(id int64) (repository.Product, error) {
	return p.model.GetProductByID(id)
}

func (p ProductService) GetAllProducts() ([]repository.Product, error) {
	return p.model.GetAllProducts()
}

func (p ProductService) UpdateProduct(id int64, name string, description string, properties map[string]string) error {
	product := repository.Product{
		Name:        name,
		Description: description,
		Properties:  properties,
	}

	return p.model.UpdateProduct(product)
}

func (p ProductService) DeleteProduct(id int64) error {
	return p.DeleteProduct(id)
}
