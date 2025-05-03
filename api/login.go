package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"media.cosasdns.com/internal"
	"media.cosasdns.com/models"
)

func HandleLogin(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	var login_params models.LoginParams

	err := json.NewDecoder(request.Body).Decode(&login_params)
	if err != nil {
		fmt.Println(err)
		internal.Log(app, "Could not parse login parameters")
		http.Error(writter, "Invalid parameters", http.StatusBadRequest)
		return
	}
	Username := login_params.Username
	Password := login_params.Password

	if Username == "" || Password == "" {
		http.Error(writter, "Credentials cannot be empty", http.StatusBadRequest)
		return
	}
	db, error_bool := internal.DB(app)
	if error_bool {
		writter.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT USER FROM USERS WHERE USER = ? AND PASSWORD = ?"

	var User string
	err = db.QueryRow(query, Username, internal.Hash(Password)).Scan(&User)
	if err != nil {
		internal.Log(app, "Login failed, no user found'")
		writter.WriteHeader(http.StatusBadRequest)
		http.Error(writter, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token := internal.RefreshToken(app, User)
	if token == "" {
		writter.WriteHeader(http.StatusInternalServerError)
		// internal.ErrorText(app, writter, "Could not generate a new token")
		http.Error(writter, "Could not generate a new token", http.StatusBadRequest)
		return
	}

	result := struct {
		Token string `json:"token"`
	}{Token: token}
	internal.WriteJsonToClient(result, writter, app)
}

func AuthMiddleware(app *models.Application, callback func(http.ResponseWriter, *http.Request, *models.Application, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO Check if the given user has permissions over the requested domain
		token := internal.GetBearerToken(r)
		if token == "" {
			internal.Log(app, "Login failed, no user found'")
			w.WriteHeader(http.StatusBadRequest)
		}
		db, err := internal.DB(app)
		if err {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer db.Close()
		query := "SELECT USER FROM USERSTOKENS WHERE TOKEN = ? AND EXPIRE_AT > CURRENT_TIMESTAMP LIMIT 1"
		var db_user string
		token_error := db.QueryRow(query, token).Scan(&db_user)
		if token_error != nil {
			internal.Log(app, "Invalid token")
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		callback(w, r, app, db_user)
	}
}
