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
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
	// TODO Create home page
	fmt.Fprintln(writter, "Home")
}
