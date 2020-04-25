package api

import (
	"bytes"
	"encoding/json"
	"net/url"
	"path"
)

type apiRequest struct {
	Method    string     `json:"method"`
	GIDList   [][]string `json:"gidlist"`
	Namespace int        `json:"namespace"`
}

func newAPIRequest(galleryURL url.URL) (*bytes.Buffer, error) {
	cleanPath := path.Clean(galleryURL.Path)
	dir, galleryToken := path.Split(cleanPath)
	galleryID := path.Base(dir)

	b := new(bytes.Buffer)
	r := apiRequest{"gdata", [][]string{{galleryID, galleryToken}}, 1}
	err := json.NewEncoder(b).Encode(r)

	return b, err
}
