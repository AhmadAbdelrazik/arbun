package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

func (app *Application) postProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Properties  map[string]string `json:"properties"`
		Vendor      string            `json:"vendor"`
		Amount      int               `json:"amount"`
		Price       float32           `json:"price"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p := domain.Product{
		Name:            input.Name,
		Description:     input.Description,
		Properties:      input.Properties,
		Vendor:          input.Vendor,
		Price:           input.Price,
		AvailableAmount: input.Amount,
	}

	v := p.Validate()
	if v != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	product, err := app.services.Products.InsertProduct(p)
	if err != nil {
		var v *validator.Validator
		switch {
		case errors.As(err, &v):
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, services.ErrDuplicateProduct):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product, err := app.services.Products.GetProductByID(id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProductNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *Application) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.services.Products.GetAllProducts()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type patchProductInput struct {
	Name            *string           `json:"name"`
	Description     *string           `json:"description"`
	Vendor          *string           `json:"vendor"`
	Properties      map[string]string `json:"properties"`
	AvailableAmount *int              `json:"available_amount"`
	Price           *float32          `json:"price"`
}

func (p patchProductInput) GenerateParams(id int64) services.UpdateProductParam {
	return services.UpdateProductParam{
		ID:              id,
		Name:            p.Name,
		Description:     p.Description,
		Vendor:          p.Vendor,
		Properties:      p.Properties,
		AvailableAmount: p.AvailableAmount,
		Price:           p.Price,
	}
}

func (app *Application) patchProduct(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input patchProductInput

	err = readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	params := input.GenerateParams(id)

	product, err := app.services.Products.UpdateProduct(params)
	if err != nil {
		var v *validator.Validator
		switch {
		case errors.As(err, &v):
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, services.ErrProductNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, services.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *Application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.Products.DeleteProduct(id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProductNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"message": "deleted successfully"}, nil)
}
