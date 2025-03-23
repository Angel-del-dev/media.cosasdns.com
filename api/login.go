package api

import (
	"fmt"
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func HandleLogin(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	err := request.ParseForm()
	if err != nil {
		internal.Log(app, "Could not parse login parameters")
		writter.WriteHeader(http.StatusBadRequest)
		return
	}

	Username := request.PostForm.Get("User")
	Password := request.PostForm.Get("Password")

	if Username == "" || Password == "" {
		writter.WriteHeader(http.StatusBadRequest)
		return
	}
	db, error_bool := internal.DB(app)
	if error_bool {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := "SELECT USER FROM USERS WHERE NAME = ? AND PASSWORD = ?"

	User := 0
	err = db.QueryRow(query, Username, internal.Hash(Password)).Scan(&User)
	if err != nil || User == 0 {
		internal.Log(app, "Login failed, no user found'")
		writter.WriteHeader(http.StatusBadRequest)
		internal.ErrorText(app, writter, "Invalid credentials")
		return
	}

	token := internal.RefreshToken(app, User)
	if token == "" {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := struct {
		Token string `json:"token"`
	}{Token: token}
	internal.WriteJsonToClient(result, writter, app)
}

func CheckToken(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	token := internal.GetBearerToken(request)

	if token == "" {
		internal.Log(app, "Token is empty")
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	origin := request.Header.Get("Origin")

	if origin != "" && origin != fmt.Sprintf("http://localhost:%d", app.Port) {
		internal.Log(app, fmt.Sprintf("Trying to check a token from a different domain '%s'", request.Header.Get("Origin")))
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	db, error_bool := internal.DB(app)
	if error_bool {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	User := 0

	query := "SELECT USER FROM USERS WHERE TOKEN = ?"
	err := db.QueryRow(query, token).Scan(&User)
	if err != nil || User == 0 {
		internal.Log(app, "Login failed, no user found'")
		writter.WriteHeader(http.StatusBadRequest)
		internal.ErrorText(app, writter, "Invalid credentials")
		return
	}

	token = internal.RefreshToken(app, User)
	if token == "" {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := struct {
		Token string `json:"token"`
	}{Token: token}
	internal.WriteJsonToClient(result, writter, app)

}
