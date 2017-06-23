package api

import (
	"encoding/json"
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	// the use for the database
	_ "github.com/lib/pq"
)

// AuthMiddleware authenticates the user logging in
func AuthMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// token is returned if the Authrization in the header matches with the users token
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})
	if err != nil {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(errMessage)
		return
	}

	if token.Valid {
		next(res, req)
	}

	if !token.Valid {
		unauthorizedMessage, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: errors.New("Unauthorized Entry").Error(),
		})
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(unauthorizedMessage)
		return
	}

}
