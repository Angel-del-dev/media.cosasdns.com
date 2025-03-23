package web

import (
	"html/template"
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

// TODO Create Login and functionality

func ServeLogin(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	tmpl, err := template.ParseFiles(
		"../templates/login/login.html",
		"../templates/generic/head.html",
		"../templates/generic/footer.html",
	)
	if err != nil {
		internal.Log(app, "Could not parse 'login' templates")
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
	err = tmpl.Execute(writter, nil)
	if err != nil {
		internal.Log(app, "Coul not serve 'login' view")
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
}
