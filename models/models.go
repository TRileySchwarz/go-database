package models

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        string `json:"email" pg:",unique"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Struct for msging service at some point
//type Msg struct {
//	ID                   uuid.UUID `json:"id" pg:",unique, type:uuid default uuid_generate_v4()"`
//	Msg string `json:"msg"`
//	Sender User `json:"sender" pg:"fk:id"`
//	Timestamp uint `json:"timestamp"`
//}

func (u *User) MarshalJSON() ([]byte, error) {
	type UserAlias User

	//fmt.Printf("\n\n This is the user values before Marshal: %v \n\n", u)

	// Alias the User struct and overloading the Password field to omitempty
	// This allows us to Unmarshal with a password field, but when Marshalling,
	// it is omited and not part of http response body
	safeUser := struct {
		// This json tag needs to be the same, or it wont have the desired affect
		Password string `json:"password,omitempty"`
		UserAlias }{
		UserAlias: UserAlias(*u),
		// Because there is no password declared in this literal it is considered empty, thus omited
	}

	//test, _ := json.Marshal(safeUser)
	//fmt.Printf("\n\n This is the user values after Marshal: %v \n\n", string(test))

	return json.Marshal(safeUser)
}

type WebTokenResponse struct {
	Token string `json:"token"`
}

type GetUserResponse struct {
	Users []User `json:"users"`
}

type PutUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type FrontEndErr struct {
	ErrorMsg string `json:"error"`
}
