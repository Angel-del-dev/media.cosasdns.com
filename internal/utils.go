package internal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"media.cosasdns.com/models"
)

func GetBearerToken(request *http.Request) string {
	return strings.Trim(strings.Replace(request.Header.Get("Authorization"), "Bearer", "", -1), " ")
}

func Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenerateToken() string {
	milli := time.Now().UnixNano() / int64(time.Millisecond)
	return Hash(fmt.Sprintf("%d", milli))
}

func ParseParams[T any](data *T, request *http.Request) {
	json.NewDecoder(request.Body).Decode(data)
}

func ErrorText(app *models.Application, writter http.ResponseWriter, message string) {
	writter.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		Log(app, "Could not generate error")
	}
	writter.Write(jsonResp)
}

func CheckMethod(writter http.ResponseWriter, request *http.Request, desired_method string) bool {
	return request.Method == desired_method
}

func WriteJsonToClient[T any](data T, writter http.ResponseWriter, app *models.Application) {
	writter.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["data"] = string(ToJson(data, writter, app))
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		Log(app, "Could not generate error")
		return
	}
	writter.Write(jsonResp)
}

func ToJson[T any](data T, writter http.ResponseWriter, app *models.Application) []byte {
	json_bytes, err := json.Marshal(data)
	if err != nil {
		Log(app, "Could not convert result to JSON")
		return []byte{}
	}
	return json_bytes
}

func FromJson[T any](buffer []byte, to_data *T, app *models.Application) {
	err := json.Unmarshal(buffer, to_data)
	if err != nil {
		Log(app, "Could not generate JSON from result")
	}
}
