package handlers

import "net/http"

func (app *Application) PostProduct(w http.ResponseWriter, r *http.Request) {
	app.Services.Products.InsertProduct()
}

func (app *Application) GetProduct(w http.ResponseWriter, r *http.Request) {

	app.Services.Products.GetProductByID()
}
func (app *Application) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	app.Services.Products.GetAllProducts()
}
func (app *Application) PatchProduct(w http.ResponseWriter, r *http.Request) {

	app.Services.Products.UpdateProduct()
}
func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	app.Services.Products.DeleteProduct()
}
