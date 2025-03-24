package api

import (
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func GetUserApplications(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	token := internal.GetBearerToken(request)
	if token == "" {
		writter.WriteHeader(http.StatusUnauthorized)
		return
	}
	db, error_bool := internal.DB(app)
	if error_bool {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := "SELECT USER FROM USERS WHERE TOKEN = ?"

	User := 0
	err := db.QueryRow(query, token).Scan(&User)
	if err != nil || User == 0 {
		internal.Log(app, "Could not obtain user from token")
		writter.WriteHeader(http.StatusUnauthorized)
		internal.ErrorText(app, writter, "Invalid credentials")
	}
	query = "SELECT DISTINCT UR.APPLICATION_NAME FROM USERROLES UR LEFT JOIN USERS U ON UR.USER = U.USER WHERE U.TOKEN = ? AND UR.PERMISSION_READ = 1 ORDER BY UR.APPLICATION_NAME ASC"

	cursor, err := db.Query(query, token)
	if err != nil {
		internal.Log(app, "Login failed, no user found'")
		writter.WriteHeader(http.StatusBadRequest)
		internal.ErrorText(app, writter, "Invalid credentials")
		return
	}

	Application := ""
	Applications := []string{}
	for cursor.Next() {
		err = cursor.Scan(&Application)
		if err != nil {
			internal.Log(app, "Could not read 'APPLICATION_NAME'")
			continue
		}
		Applications = append(Applications, Application)
	}

	token = internal.RefreshToken(app, User)
	result := struct {
		Applications []string `json:"applications"`
		Token        string   `json:"token"`
	}{Applications: Applications, Token: token}
	internal.WriteJsonToClient(result, writter, app)
}

func CreateApplication(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	token := internal.GetBearerToken(request)
	if token == "" {
		writter.WriteHeader(http.StatusUnauthorized)
		return
	}
	db, error_bool := internal.DB(app)
	if error_bool {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := "SELECT USER FROM USERS WHERE TOKEN = ?"

	User := 0
	err := db.QueryRow(query, token).Scan(&User)
	if err != nil || User == 0 {
		internal.Log(app, "Could not obtain user from token")
		writter.WriteHeader(http.StatusUnauthorized)
		internal.ErrorText(app, writter, "Invalid credentials")
	}

	token = internal.RefreshToken(app, User)

	stmt, err := db.Prepare("INSERT INTO APPLICATIONS (NAME) VALUES (?)")
	if err != nil {
		internal.Log(app, "Could not create prepared statement for APPLICATION creation")
	}

	_, err = stmt.Exec(token, internal.GenerateRandomString(10))
	if err != nil {
		internal.Log(app, "Could not execute APPLICATION creation")
	}

	// TODO assign a set name to the application && create a link between a domain, an application and an user

	result := struct {
		Token string `json:"token"`
	}{Token: token}
	internal.WriteJsonToClient(result, writter, app)
}
