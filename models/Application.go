package models

import (
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Port     int64
	Log      bool
	LogRoute string
}

func (app *Application) Init() {
	app.Log = true
	app.Port = 80
	app.LogRoute = "../logs"
	app.parseArguments()
}

func (app *Application) parseArguments() {
	Args := os.Args[1:]
	for _, value := range Args {
		arg := strings.Split(value, "=")
		if len(arg) < 2 {
			continue
		}

		switch strings.ToUpper(arg[0]) {
		case "PORT":
			num, err := strconv.ParseInt(arg[1], 10, 64)
			if err != nil {
				continue
			}
			app.Port = num
		case "LOG":
			if arg[1] == "false" {
				app.Log = false
			}
		case "LOGROUTE":
			app.LogRoute = arg[1]
		}
	}
}
