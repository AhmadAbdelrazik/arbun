package handlers

import "net/http"

func (app *Application) PostProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Properties  map[string]string `json:"properties"`
	}

	id, err := app.services.Products.InsertProduct(input.Name, input.Description, input.Properties)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"id": id}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetProduct(w http.ResponseWriter, r *http.Request) {

	app.services.Products.GetProductByID()
}
func (app *Application) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	app.services.Products.GetAllProducts()
}
func (app *Application) PatchProduct(w http.ResponseWriter, r *http.Request) {

	app.services.Products.UpdateProduct()
}
func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	app.services.Products.DeleteProduct()
}
