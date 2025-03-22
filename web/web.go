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

	http.HandleFunc("/file/{application}/{resource}", func(writter http.ResponseWriter, request *http.Request) { api.GetResource(writter, request, app) })
	http.HandleFunc("/favicon.ico", doNothing)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.Port), nil)
	if err != nil {
		internal.Log(app, "Could not start the server")
	}
}
