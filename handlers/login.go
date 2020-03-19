package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("pass")

	row := db.QueryRow("select id, name, info from tes_bd.users where login = ? and pass = ?", login, password)

	p := formatData{}
	dataDB := row.Scan(&p.id, &p.Name, &p.Info)

	if dataDB == sql.ErrNoRows {
		data := formatData{
			Message: "Oh, you were mistaken in the password or login ... oh my God what to do now =(",
		}
		tempLogin.Execute(w, data)
	} else if dataDB != nil {
		data := formatData{
			Message: "Oh no, something happened to the database",
		}
		tempLogin.Execute(w, data)

	} else {
		value := map[string]string{
			"login": login,
			"name":  p.Name,
			"info":  p.Info,
		}
		if encoded, err := securecookie.EncodeMulti("auth", value, cookies["current"]); err == nil {
			expire := time.Now().Local().AddDate(1, 0, 0)
			cookie := &http.Cookie{
				Name:    "auth",
				Value:   encoded,
				Path:    "/",
				Expires: expire,
			}
			http.SetCookie(w, cookie)
		}
		http.ServeFile(w, r, "/")
	}
}
