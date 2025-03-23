package web

import (
	"fmt"
	"net/http"

	"media.cosasdns.com/models"
)

func Panel(writter http.ResponseWriter, request *http.Request, app *models.Application) {
	fmt.Println("Panel")
}
