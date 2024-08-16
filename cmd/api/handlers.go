package main

import (
	"net/http"
	"strconv"
	"strings"

	"miruchigawa.moe/restapi/internal/funcs/anime"
	"miruchigawa.moe/restapi/internal/funcs/manga"
	"miruchigawa.moe/restapi/internal/funcs/downloader"
	"miruchigawa.moe/restapi/internal/response"
	"miruchigawa.moe/restapi/internal/validator"
	"miruchigawa.moe/restapi/internal/version"
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
	} else {
		v.AddError("query can't be empty!")
	}

	if pageQuery := query.Get("page"); pageQuery != "" {
		if num, err := strconv.Atoi(pageQuery); err == nil || num < 1 {
			page = num
		} else {
			v.AddError("Invalid page number format!")
		}
	} else {
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
		"Status":  "OK",
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
	} else {
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
		"Status":  "OK",
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
	} else {
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
		"Status":  "OK",
		"Message": result,
	}

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) mangaSearch(w http.ResponseWriter, r *http.Request) {
	var name string
	var page int
	var limit int
	query := r.URL.Query()
	v := validator.Validator{}

	if queryName := query.Get("query"); queryName != "" {
		name = strings.TrimSpace(queryName)
		v.Check(len(name) > 0, "query can't be empty!")
	} else {
		v.AddError("query can't be empty!")
	}

	if pageQuery := query.Get("page"); pageQuery != "" {
		if num, err := strconv.Atoi(pageQuery); err == nil || num < 1 {
			page = num
		} else {
			v.AddError("Invalid page number format!")
		}
	} else {
		page = 1
	}

	if limitQuery := query.Get("page"); limitQuery != "" {
		if num, err := strconv.Atoi(limitQuery); err == nil || num < 1 {
			limit = num
		} else {
			v.AddError("Invalid limit number format!")
		}
	} else {
		limit = 20
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := manga.Search(name, page, limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status":  "OK",
		"Message": result,
	}

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) mediafire(w http.ResponseWriter, r *http.Request) {
	var url string
	query := r.URL.Query()
	v := validator.Validator{}

	if queryUrl := query.Get("url"); queryUrl != "" {
		url = strings.TrimSpace(queryUrl)
		v.Check(len(url) > 0, "url can't be empty!")
	} else {
		v.AddError("url can't be empty!")
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := downloader.GetMediafireInfo(url)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status":  "OK",
		"Message": result,
	}

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) tiktokDownloader(w http.ResponseWriter, r *http.Request) {
	var url string
	query := r.URL.Query()
	v := validator.Validator{}

	if queryUrl := query.Get("url"); queryUrl != "" {
		url = strings.TrimSpace(queryUrl)
		v.Check(len(url) > 0, "url can't be empty!")
	} else {
		v.AddError("url can't be empty!")
	}

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := downloader.TiktokDownloader(url)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]any{
		"Status":  "OK",
		"Message": result,
	}

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}

}
