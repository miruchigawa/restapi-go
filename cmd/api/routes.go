package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.NotFoundHandler = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)

	mux.HandleFunc("/status", app.status).Methods("GET")
	mux.HandleFunc("/anime/search", app.animeSearch).Methods("GET")
	mux.HandleFunc("/anime/info", app.animeInfo).Methods("GET")
	mux.HandleFunc("/manga/search", app.mangaSearch).Methods("GET")
	mux.HandleFunc("/downloader/mediafire", app.mediafire).Methods("GET")
	mux.HandleFunc("/downloader/tiktok", app.tiktokDownloader).Methods("GET")

	return mux
}
