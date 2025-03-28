package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"media.cosasdns.com/models"
	"media.cosasdns.com/pkg"
)

func DB(app *models.Application) (*sql.DB, bool) {
	configuration, err := pkg.GetIni()

	if err {
		return nil, true
	}

	section := configuration.Section("database")

	host := section.Key("DB_HOST").String()
	port := section.Key("DB_PORT").String()
	user := section.Key("DB_USER").String()
	password := section.Key("DB_PASSWORD").String()
	dbname := section.Key("DB_NAME").String()

	if host == "" || user == "" || password == "" || dbname == "" {
		Log(app, "The ini file is not properly configured")
	}

	connection_query := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	db, conn_error := sql.Open("mysql", connection_query)

	if conn_error != nil {
		Log(app, "Could not connect to database")
		return nil, true
	}
	return db, false
}

func RefreshToken(app *models.Application, User int) string {
	token := GenerateToken()

	db, error_bool := DB(app)
	if error_bool {
		return ""
	}

	stmt, err := db.Prepare("UPDATE USERS SET TOKEN = ? WHERE USER = ?")
	if err != nil {
		Log(app, "Could not create prepared statement for token storage")
	}

	_, err = stmt.Exec(token, User)
	if err != nil {
		Log(app, "Could not execute token storage")
	}

	return token
}
