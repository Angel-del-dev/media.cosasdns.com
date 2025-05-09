package web

import (
	"fmt"
	"net/http"

	"media.cosasdns.com/api"
	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func doNothing(rw http.ResponseWriter, req *http.Request) {}

func InitWebHandler(app *models.Application) {
	internal.Log(app, "Starting http server")
	internal.Log(app, fmt.Sprintf("Started on port :%d", app.Port))

	mux := http.NewServeMux()

	// Api
	mux.HandleFunc("/api/login", func(writter http.ResponseWriter, request *http.Request) { api.HandleLogin(writter, request, app) })

	mux.HandleFunc("/file/{resource}", func(writter http.ResponseWriter, request *http.Request) { api.GetResource(writter, request, app) })
	mux.HandleFunc("/file/add-from-string", api.AuthMiddleware(app, api.Handle))
	mux.HandleFunc("/file/all", api.AuthMiddleware(app, api.GetAll))
	mux.HandleFunc("/file/all-from-type/{type}", api.AuthMiddleware(app, api.GetAllFromMime))
	mux.HandleFunc("/file/unlink/{resource}", api.AuthMiddleware(app, api.RemoveFile))

	mux.HandleFunc("/favicon.ico", doNothing)
	mux.HandleFunc("/login", func(writter http.ResponseWriter, request *http.Request) { Login(writter, request, app) })
	mux.HandleFunc("/", func(writter http.ResponseWriter, request *http.Request) { ServeHome(writter, request, app) })

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		internal.Log(app, "Could not start the server")
	}
}
