package models

type ResourceParams struct {
	FileString string `json:"filestring"`
	MimeType   string `json:"mimetype"`
}

type ResourceCatalog struct {
	Token     string   `json:"token"`
	Resources []string `json:"resources"`
}
