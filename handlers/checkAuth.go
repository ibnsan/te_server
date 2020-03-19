package handlers

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

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
