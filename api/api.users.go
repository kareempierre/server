package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UserCredentials requires an Email and Password
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserConstruct is what you get back from the database
type UserConstruct struct {
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Organization string
	Admin        bool
	Creator      bool
	register     time.Time
}

// UsersList is a struct for the information returned on all users
type UsersList struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}

var (
	// DBUser is the current logged in user
	DBUser UserConstruct
)

// CreateUser creates a new user
func CreateUser(res http.ResponseWriter, req *http.Request) {
	// Temporary user variable when grabbing the information from the body
	var user UserConstruct

	// Decode the information and pass into the temporary variable
	err := json.NewDecoder(req.Body).Decode(&user)
	if err, ok := OnError(err, http.StatusBadRequest); !ok {
		res.Write(err)
		return
	}

	//Hash the current password from the temporary variable
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err, ok := OnError(err, http.StatusBadRequest); !ok {
		res.Write(err)
		return
	}

	// Prepare Qeury returns statment to insert new user
	stmt, err := DB.Prepare(`INSERT INTO users(firstname, lastname, email, admin, creator, organization, password, register)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8);`)
	defer stmt.Close()
	if err, ok := OnError(err, http.StatusInternalServerError); !ok {
		res.Write(err)
		return
	}

	// Execute the prepared statement and add to the DB
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, false, false, user.Organization, string(hashedPassword), time.Now())
	if err, ok := OnError(err, http.StatusInternalServerError); !ok {
		res.Write(err)
		return
	}

	// Generate a new signing key to assign to the created user
	token := jwt.New(jwt.SigningMethodRS256)
	SignedKey, err := token.SignedString(SignKey)
	if err, ok := OnError(err, http.StatusUnauthorized); !ok {
		res.Write(err)
		return
	}

	// Create the new user and Marshal data to be sent back
	userInfo, err := json.Marshal(struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Token     string `json:"token"`
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     SignedKey,
	})
	if err, ok := OnError(err, http.StatusInternalServerError); !ok {
		res.Write(err)
		return
	}

	// Respond with user information and json key
	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// AuthHandler is used to authenticate the user logging in
func AuthHandler(res http.ResponseWriter, req *http.Request) {
	// Temporary login variable
	var user UserCredentials

	// Get information from body of the request
	err := json.NewDecoder(req.Body).Decode(&user)
	if err, ok := OnError(err, http.StatusBadRequest); !ok {
		res.Write(err)
		return
	}

	// Query row from database using the information gained from the temporary variable
	row := DB.QueryRow(`SELECT firstname, lastname, email, organization, creator, admin, password FROM users WHERE email=$1`, user.Email)
	rowErr := row.Scan(&DBUser.FirstName, &DBUser.LastName, &DBUser.Email, &DBUser.Organization,
		&DBUser.Creator, &DBUser.Admin, &DBUser.Password)

	if rowErr != nil || rowErr == sql.ErrNoRows {
		errMessage, _ := json.Marshal(ErrorResponse{
			Type:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(errMessage)
		return
	}

	// Check password to make sure that the using logging in is the correct user
	err = bcrypt.CompareHashAndPassword([]byte(DBUser.Password), []byte(user.Password))
	if err, ok := OnError(err, http.StatusUnauthorized); !ok {
		res.Write(err)
		return
	}

	// Generate token for the user logging in
	token := jwt.New(jwt.SigningMethodRS256)

	SignedKey, err := token.SignedString(SignKey)
	if err, ok := OnError(err, http.StatusBadRequest); !ok {
		res.Write(err)
		return
	}

	// Create user information collected from the DB
	userInfo, err := json.Marshal(struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Token     string `json:"token"`
	}{
		FirstName: DBUser.FirstName,
		LastName:  DBUser.LastName,
		Token:     SignedKey,
	})
	if err, ok := OnError(err, http.StatusInternalServerError); !ok {
		res.Write(err)
		return
	}

	// Respond with information
	res.Header().Set("Content-Type", "application/json")
	res.Write(userInfo)

}

// UsersHandler returns all users for a specific organization
func UsersHandler(res http.ResponseWriter, req *http.Request) {
	// user list variable
	var usersList UsersList
	// user array of user lists
	var usersArray []UsersList

	// Get all users listed in a specific organization
	rows, err := DB.Query(`SELECT firstname, lastname, email, admin FROM users WHERE organization=$1;`, DBUser.Organization)
	if DBUser.Organization == "all" && DBUser.Creator == true {
		rows, err = DB.Query(`SELECT firstname, lastname, email, admin FROM users;`)
	}
	if err, ok := OnError(err, http.StatusInternalServerError); !ok {
		res.Write(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&usersList.FirstName, &usersList.LastName, &usersList.Email, &usersList.Admin)
		if err, ok := OnError(err, http.StatusInternalServerError); !ok {
			res.Write(err)
			return
		}
		usersArray = append(usersArray, usersList)
	}

	// If user requesting this information is an admin return all the users for that organization
	// else The user requesting this information is returned a bool specifying they are not an admin
	if DBUser.Admin == true {
		users, err := json.Marshal(usersArray)
		if err, ok := OnError(err, http.StatusInternalServerError); !ok {
			res.Write(err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(users)
	} else {
		unAuthUser, err := json.Marshal(struct {
			Type  int  `json:"type"`
			Admin bool `json:"admin"`
		}{
			Type:  http.StatusUnauthorized,
			Admin: false,
		})
		if err, ok := OnError(err, http.StatusUnauthorized); !ok {
			res.Write(err)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(unAuthUser)
	}
}

// ViewUserHandler requests information based on the user being viewed in the admin section
func ViewUserHandler(res http.ResponseWriter, req *http.Request) {
	//var email string

	//row := DB.QueryRow(`SELECT firstname, lastname, email`)
}
