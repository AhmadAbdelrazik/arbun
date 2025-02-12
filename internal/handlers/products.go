package handlers

import (
	"AhmadAbdelrazik/arbun/internal/services"
	"net/http"
)

func (app *Application) PostProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Properties  map[string]string `json:"properties"`
		Vendor      string            `json:"vendor"`
		Amount      int               `json:"amount"`
	}
	// TODO: Handle input validation

	params := services.InsertProductParam{
		Name:            input.Name,
		Description:     input.Description,
		Properties:      input.Properties,
		Vendor:          input.Vendor,
		AvailableAmount: input.Amount,
	}
	id, err := app.services.Products.InsertProduct(params)
	if err != nil {
		// TODO: Handle different errors
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"id": id}, nil)
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
		// TODO: Handle different errors
		app.badRequestResponse(w, r, err)
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
		// TODO: Handle different errors
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *Application) PatchProduct(w http.ResponseWriter, r *http.Request) {

}
func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {
}
