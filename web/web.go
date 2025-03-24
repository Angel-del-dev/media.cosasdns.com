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

	fs := http.FileServer(http.Dir("../static"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	// Api calls
	mux.HandleFunc("/check-token", func(writter http.ResponseWriter, request *http.Request) { api.CheckToken(writter, request, app) })
	mux.HandleFunc("/get-user-applications", func(writter http.ResponseWriter, request *http.Request) {
		api.GetUserApplications(writter, request, app)
	})

	mux.HandleFunc("/create-application", func(writter http.ResponseWriter, request *http.Request) {
		api.CreateApplication(writter, request, app)
	})
	// Other calls

	mux.HandleFunc("/file/{application}/{resource}", func(writter http.ResponseWriter, request *http.Request) { api.GetResource(writter, request, app) })
	mux.HandleFunc("/favicon.ico", doNothing)
	mux.HandleFunc("/login", func(writter http.ResponseWriter, request *http.Request) { Login(writter, request, app) })
	mux.HandleFunc("/panel", func(writter http.ResponseWriter, request *http.Request) { Panel(writter, request, app) })
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
