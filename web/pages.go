package web

import (
	"fmt"
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func ServeHome(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	if request.URL.Path != "/" {
		internal.Log(app, fmt.Sprintf("Invalid route '%s'", request.URL.Path))
		if internal.CheckMethod(writter, request, "GET") {
			http.Redirect(writter, request, "/", http.StatusSeeOther)
		} else {
			writter.WriteHeader(http.StatusNotFound)
		}
		return
	}
	// Redirect to the api reference
	http.Redirect(writter, request, "http://docs.angelnovo.es", http.StatusSeeOther)
}
