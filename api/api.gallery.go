package api

import (
	"fmt"
	"net/http"
)

// ViewGalleryHandler gets the specific images to be used by the client
func ViewGalleryHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Gallery endpoint hit")
}

// UploadImageHandler posts uploaded
func UploadImageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Uploaded Image endpoint hit")
}

// GetImageHandler gets a single Image
func GetImageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("GetImage Handler endpoint hit")
}

// TotalByteHandler displays the total amount of space used for a specific organization
func TotalByteHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Total bytes used endpoint hit")
}
