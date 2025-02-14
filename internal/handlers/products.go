package handlers

import (
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

type postProductInput struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Properties  map[string]string `json:"properties"`
	Vendor      string            `json:"vendor"`
	Amount      int               `json:"amount"`
}

func (p postProductInput) GenerateParams() services.InsertProductParam {
	return services.InsertProductParam{
		Name:            p.Name,
		Description:     p.Description,
		Properties:      p.Properties,
		Vendor:          p.Vendor,
		AvailableAmount: p.Amount,
	}
}

func (app *Application) PostProduct(w http.ResponseWriter, r *http.Request) {
	var input postProductInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	params := input.GenerateParams()
	product, err := app.services.Products.InsertProduct(params)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDuplicateProduct):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *Application) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.services.Products.GetAllProducts()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
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
}

func (p patchProductInput) GenerateParams(id int64) services.UpdateProductParam {
	return services.UpdateProductParam{
		ID:              id,
		Name:            p.Name,
		Description:     p.Description,
		Vendor:          p.Vendor,
		Properties:      p.Properties,
		AvailableAmount: p.AvailableAmount,
	}
}

func (app *Application) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input patchProductInput

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	params := input.GenerateParams(id)

	product, err := app.services.Products.UpdateProduct(params)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProductNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, services.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "deleted successfully"}, nil)
}
