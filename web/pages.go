package web

import (
	"fmt"
	"net/http"
)

func ServeHome(writter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Redirect(writter, request, "/", http.StatusSeeOther)
		return
	}
	// TODO Create home page
	fmt.Fprintln(writter, "Home")
}
