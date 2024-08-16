package main

import (
	"strings"
	"strconv"
	"net/http"

	"miruchigawa.moe/restapi/internal/response"
	"miruchigawa.moe/restapi/internal/version"
	"miruchigawa.moe/restapi/internal/funcs/anime"
	"miruchigawa.moe/restapi/internal/validator"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Status": "OK",
		"Message": map[string]any{
			"Version": version.Get(),
		},
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) animeSearch(w http.ResponseWriter, r *http.Request) {
	var name string
	var page int
	query := r.URL.Query()
	v := validator.Validator{}

	if queryName := query.Get("query"); queryName != "" {
		name = strings.TrimSpace(queryName)
		v.Check(len(name) > 0, "query can't be empty!")
	}else{
		v.AddError("query can't be empty!")
	}

	if pageQuery := query.Get("page"); pageQuery != "" {
		if num, err := strconv.Atoi(pageQuery); err == nil || num < 1 {
			page = num
		}else{
			v.AddError("Invalid page number format!")
		}
	}else{
		page = 1
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}
	
	result, err := anime.Search(name, page)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status": "OK",
		"Message": result,
	}	

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) animeInfo(w http.ResponseWriter, r *http.Request) {
	var id string
	query := r.URL.Query()
	v := validator.Validator{}

	if queryId := query.Get("id"); queryId != "" {
		id = strings.TrimSpace(queryId)
		v.Check(len(id) > 0, "id can't be empty!")
	}else{
		v.AddError("id can't be empty!")
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := anime.Info(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status": "OK",
		"Message": result,
	}	

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) animeDownload(w http.ResponseWriter, r *http.Request) {
	var id string
	query := r.URL.Query()
	v := validator.Validator{}

	if queryId := query.Get("id"); queryId != "" {
		id = strings.TrimSpace(queryId)
		v.Check(len(id) > 0, "id can't be empty!")
	}else{
		v.AddError("id can't be empty!")
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := anime.Downloads(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status": "OK",
		"Message": result,
	}	

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}	
}
