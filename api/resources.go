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

	query := "SELECT R.RESOURCE FROM RESOURCESDOMAINS RD JOIN RESOURCES R ON RD.RESOURCE = R.RESOURCE WHERE RD.DOMAIN = ? AND RD.RESOURCE = ?"

	file_name := ""

	query_error := db.QueryRow(query, request.Header.Get("Origin"), request.PathValue("resource")).Scan(&file_name)
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
	writter.Header().Set("Content-Type", file_name)
	writter.Write(file_bytes)
}
