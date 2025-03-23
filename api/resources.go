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
	// Add method control and log
	db, err := internal.DB(app)
	if err {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT TYPE FROM RESOURCES WHERE APPLICATION_NAME = ?"
	// If Origin is empty, it means the request is from the same origin as media.cosasdns.com
	if request.Header.Get("Origin") != "" {
		query += fmt.Sprintf(" AND DOMAIN = '%s'  ", request.Header.Get("Origin"))
	}
	query += " AND NAME = ? "

	file_type := ""

	query_error := db.QueryRow(query, request.PathValue("application"), request.PathValue("resource")).Scan(&file_type)
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error obtaining resource '%s'", request.URL.Path))
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	file_route := fmt.Sprintf("../files/%s/%s", request.PathValue("application"), request.PathValue("resource"))
	file_bytes, file_error := os.ReadFile(file_route)

	if file_error != nil {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	writter.WriteHeader(http.StatusFound)
	writter.Header().Set("Content-Type", file_type)
	writter.Write(file_bytes)
}
