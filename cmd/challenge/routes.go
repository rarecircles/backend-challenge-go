package main

import "github.com/julienschmidt/httprouter"

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/tokens", app.tokens)
	return router
}
