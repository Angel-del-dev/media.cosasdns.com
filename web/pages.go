package web

import (
	"fmt"
	"net/http"
)

func Root(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Route not found")
}
