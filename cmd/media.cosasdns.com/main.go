package main

import (
	"media.cosasdns.com/models"
	"media.cosasdns.com/web"
)

func main() {
	app := models.Application{}
	app.Init()
	web.InitWebHandler(&app)
}
