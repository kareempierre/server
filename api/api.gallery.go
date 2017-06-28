package api

import (
	"fmt"
	"net/http"
)

// GalleryHandler gets the specific images to be used by the client
func GalleryHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Gallery endpoint hit")
}
