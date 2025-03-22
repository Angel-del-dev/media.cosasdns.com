package pkg

import (
	"gopkg.in/ini.v1"
)

func GetIni() (*ini.File, bool) {
	inidata, err := ini.Load("./application.ini")
	if err != nil {
		return nil, true
	}
	return inidata, false
}
