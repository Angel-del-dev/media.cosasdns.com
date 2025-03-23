package web

import (
	"net/http"

	"media.cosasdns.com/api"
	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

// TODO Create Login and functionality

func ServeLogin(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	internal.ExecuteTemplates(writter, request, app, "../templates/login/login.html", "Login")
}

func Login(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	if request.Method == "GET" {
		ServeLogin(writter, request, app)
	} else {
		api.HandleLogin(writter, request, app)
	}
}
