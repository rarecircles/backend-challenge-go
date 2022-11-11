package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type response struct {
	Tokens []tokenResponse `json:"tokens"`
}

func (app *application) tokens(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	re := response{}

	queryParams := r.URL.Query()
	q := queryParams.Get("q")
	if q == "" {
		app.logger.Error("empty q parameter")
		app.writeJSON(w, 200, re, nil)
		return
	}

	tR := getData(app, q)
	re.Tokens = tR
	app.writeJSON(w, 200, re, nil)
	return
}

type tokenAddress struct {
	Address string `json:"address"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
