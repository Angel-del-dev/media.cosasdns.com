package internal

import (
	"html/template"
	"net/http"

	"media.cosasdns.com/models"
)

func ExecuteTemplates(
	writter http.ResponseWriter,
	request *http.Request,
	app *models.Application,
	main_template_route string,
	title string,
) {
	tmpl, err := template.ParseFiles(
		main_template_route,
		"../templates/generic/head.html",
		"../templates/generic/footer.html",
	)
	if err != nil {
		Log(app, "Could not parse 'login' templates")
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
	err = tmpl.Execute(writter, models.Template{Title: title})
	if err != nil {
		Log(app, "Coul not serve 'login' view")
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
}
