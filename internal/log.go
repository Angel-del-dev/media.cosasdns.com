package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	"media.cosasdns.com/models"
)

func Log(app *models.Application, message string) {
	if !app.Log {
		return
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/log.txt", app.LogRoute), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s %s %s\n", time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"), message)); err != nil {
		log.Println(err)
	}
}
