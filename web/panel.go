package web

import (
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func Panel(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	internal.ExecuteTemplates(writter, request, app, "../templates/panel/panel.html", "Panel")
}
