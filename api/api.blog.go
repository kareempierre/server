package api

import (
	"fmt"
	"net/http"
)

// BlogHandler grabs all written blogs based on the calling organization
func BlogHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Blog endpoint hit")
}
