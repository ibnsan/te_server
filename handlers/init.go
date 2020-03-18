package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"


	"github.com/gorilla/securecookie"
)

var hashKey = []byte("very-secret")
var blockKey = []byte("a-lot-secret")
var s = securecookie.New(hashKey, blockKey)

type formatData struct {
	id      int
	Login   string
	pass    string
	Message string
	Name    string
	Info    string
}

var db *sql.DB
var tempLogin, tempRegister, tempHome *template.Template

func init() {
	var err error

	db, err = sql.Open("mysql", "auser:12345678@/tes_bd")
	if err != nil {
		panic(err)
	}

	tempLogin, err = template.ParseFiles("pages/login.html")
	if err != nil {
		panic(err)
	}

	tempRegister, err = template.ParseFiles("pages/registration.html")
	if err != nil {
		panic(err)
	}

	tempHome, err = template.ParseFiles("pages/home.html")
	if err != nil {
		panic(err)
	}

}

var cookies = map[string]*securecookie.SecureCookie{
	"previous": securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	),
	"current": securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	),
}

