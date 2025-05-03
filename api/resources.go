package api

import (
	"fmt"
	"net/http"
	"os"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

// TODO Handle files with encryption/decryption
// TODO Limit application max size (Based on a plan?)

func GetResource(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	db, err := internal.DB(app)
	if err {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT R.TYPE FROM RESOURCESDOMAINS RD JOIN RESOURCES R ON RD.RESOURCE = R.RESOURCE WHERE RD.DOMAIN = ? AND RD.RESOURCE = ?"

	file_type := ""

	query_error := db.QueryRow(query, request.Header.Get("Origin"), request.PathValue("resource")).Scan(&file_type)
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error obtaining resource '%s'", request.URL.Path))
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	file_route := fmt.Sprintf("../files/%s", request.PathValue("resource"))
	file_bytes, file_error := os.ReadFile(file_route)

	if file_error != nil {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	writter.WriteHeader(http.StatusFound)
	writter.Header().Set("Content-Type", file_type)
	writter.Write(file_bytes)
}

func addResource(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	fmt.Println("Add resource")
}

func Handle(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	if internal.CheckMethod(writter, request, "POST") {
		addResource(writter, request, app)
		return
	}

	writter.WriteHeader(http.StatusNotFound)
}
