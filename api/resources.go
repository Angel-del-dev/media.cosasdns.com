package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

// TODO Handle files with encryption/decryption
// TODO Limit application max size (Based on a plan?)
// TODO Add better error returning

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

func addResource(writter http.ResponseWriter, request *http.Request, app *models.Application, user string) {
	var resource_params models.ResourceParams

	err := json.NewDecoder(request.Body).Decode(&resource_params)
	if err != nil {
		internal.Log(app, "Invalid resource parameters, please read the docs")
		http.Error(writter, "Invalid resource parameters, please read the docs", http.StatusBadRequest)
		return
	}

	if resource_params.MimeType == "" {
		internal.Log(app, "MimeType must be provided")
		http.Error(writter, "MimeType must be provided", http.StatusBadRequest)
		return
	}

	decoded, error_b64 := base64.StdEncoding.DecodeString(resource_params.FileString)
	if error_b64 != nil {
		internal.Log(app, "The provided filestring is not a valid Base 64 encoded string")
		http.Error(writter, "The provided filestring is not a valid Base 64 encoded string", http.StatusBadRequest)
		return
	}

	db, db_err := internal.DB(app)
	if db_err {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	FileName := internal.GenerateRandomString(100)

	// Save the resource to the db
	stmt, err := db.Prepare("INSERT INTO RESOURCES (RESOURCE, TYPE) VALUES (?, ?)")
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(FileName, resource_params.MimeType)
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Link the resource to the current domain
	stmt, err = db.Prepare("INSERT INTO RESOURCESDOMAINS (DOMAIN, RESOURCE) VALUES (?, ?)")
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(request.Header.Get("Origin"), FileName)
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Save the file
	current_directory, _ := os.Getwd()
	file, err := os.Create(current_directory + "/../files/" + FileName)
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(string(decoded))
	if err != nil {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := internal.RefreshToken(app, user)
	if token == "" {
		writter.WriteHeader(http.StatusInternalServerError)
		http.Error(writter, "Could not generate a new token", http.StatusBadRequest)
		return
	}
	result := struct {
		Token string `json:"token"`
	}{Token: token}
	internal.WriteJsonToClient(result, writter, app)
}

func Handle(writter http.ResponseWriter, request *http.Request, app *models.Application, user string) {
	if internal.CheckMethod(writter, request, "POST") {
		addResource(writter, request, app, user)
		return
	}

	writter.WriteHeader(http.StatusNotFound)
}
