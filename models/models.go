package models

import (
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

type UserNoPwd struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type WebTokenResponse struct {
	Token string `json:"token"`
}

type GetUserResponse struct {
	Users []UserNoPwd `json:"users"`
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
