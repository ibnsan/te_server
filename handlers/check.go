package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
)

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

func Сheckauth(w http.ResponseWriter, r *http.Request) {

	if cookie, err := r.Cookie("auth"); err == nil {
		value := make(map[string]string)
		err = securecookie.DecodeMulti("auth", cookie.Value, &value, cookies["current"], cookies["previous"])
		if err == nil {
			data := formatData{
				Name:  value["name"],
				Login: value["login"],
				Info:  value["info"],
			}
			tempHome.Execute(w, data)
		}
	} else {
		data := formatData{
			Message: "make a mistake and I will remember that о_о",
		}
		tempLogin.Execute(w, data)
	}

}
