package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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

func GetAllFromDomainAndMime(origin string, TYPE string, writter http.ResponseWriter, app *models.Application, user string) models.ResourceCatalog {
	query := " SELECT RD.RESOURCE FROM RESOURCESDOMAINS RD LEFT JOIN RESOURCES R ON RD.RESOURCE = R.RESOURCE WHERE RD.DOMAIN = ? "

	if TYPE != "" {
		query += " AND (UPPER(R.TYPE) LIKE '" + TYPE + "/%' OR UPPER(R.TYPE) LIKE '%/" + TYPE + "') "
	}

	db, err := internal.DB(app)
	if err {
		writter.WriteHeader(http.StatusInternalServerError)
		return models.ResourceCatalog{}
	}

	defer db.Close()

	rows, query_error := db.Query(query, origin)
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error obtaining resources from origin '%s'", origin))
		writter.WriteHeader(http.StatusInternalServerError)
		return models.ResourceCatalog{}
	}

	Resources := []string{}

	for rows.Next() {
		var resource_col string
		rows.Scan(&resource_col)
		Resources = append(Resources, resource_col)
	}

	token := internal.RefreshToken(app, user)
	if token == "" {
		writter.WriteHeader(http.StatusInternalServerError)
		http.Error(writter, "Could not generate a new token", http.StatusBadRequest)
		return models.ResourceCatalog{}
	}

	return models.ResourceCatalog{Token: token, Resources: Resources}
}

func GetAll(writter http.ResponseWriter, request *http.Request, app *models.Application, user string) {
	if !internal.CheckMethod(writter, request, "GET") {
		writter.WriteHeader(http.StatusNotFound)
		return
	}
	result := GetAllFromDomainAndMime(request.Header.Get("Origin"), "", writter, app, user)
	internal.WriteJsonToClient(result, writter, app)
}

func GetAllFromMime(writter http.ResponseWriter, request *http.Request, app *models.Application, user string) {
	if !internal.CheckMethod(writter, request, "GET") {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	result := GetAllFromDomainAndMime(request.Header.Get("Origin"), strings.ToUpper(request.PathValue("type")), writter, app, user)
	internal.WriteJsonToClient(result, writter, app)
}

func RemoveFile(writter http.ResponseWriter, request *http.Request, app *models.Application, user string) {
	if !internal.CheckMethod(writter, request, "DELETE") {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	db, err := internal.DB(app)
	if err {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "DELETE FROM RESOURCESDOMAINS WHERE DOMAIN = ? AND RESOURCE = ?"

	// Removal of the link -> resource/domain

	_, query_error := db.Exec(query, request.Header.Get("Origin"), request.PathValue("resource"))
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error unlinking resource '%s'", request.PathValue("resource")))
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get count of the domains with the resource
	CountRemainingDomains := 0
	query = "SELECT COUNT(1) FROM RESOURCESDOMAINS WHERE DOMAIN = ? AND RESOURCE = ?"
	query_error = db.QueryRow(query, request.Header.Get("Origin"), request.PathValue("resource")).Scan(&CountRemainingDomains)
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error counting resource '%s' domains", request.PathValue("resource")))
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	if CountRemainingDomains > 0 {
		writter.WriteHeader(http.StatusOK)
		return
	}
	// Remove File / Delete resource

	query = "DELETE FROM RESOURCES WHERE RESOURCE = ?"
	_, query_error = db.Exec(query, request.PathValue("resource"))
	if query_error != nil {
		internal.Log(app, fmt.Sprintf("Error removing resource '%s'", request.PathValue("resource")))
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("../files/%s", request.PathValue("resource"))
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err := os.Remove(filePath)
		if err != nil {
			internal.Log(app, fmt.Sprintf("Error removing file '%s'", request.PathValue("resource")))
			writter.WriteHeader(http.StatusInternalServerError)
		}
	}

	writter.WriteHeader(http.StatusOK)
}
