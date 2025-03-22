package web

import (
	"fmt"
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func InitWebHandler(app *models.Application) {
	internal.Log(app, "Starting http server")
	internal.Log(app, fmt.Sprintf("Started on port :%d", app.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", app.Port), nil)
}
