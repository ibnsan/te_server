package handlers

import (
	"net/http"
	// "github.com/ibnsan/te_server/handlers"
)

// test
func Сheckauth(w http.ResponseWriter, r *http.Request) {

	data := formatData{
		Message: "make a mistake and I will remember that о_о",
	}
	tempLogin.Execute(w, data)

	// if cookie, err := r.Cookie("auth"); err == nil {
	// 	value := make(map[string]string)
	// 	err = securecookie.DecodeMulti("auth", cookie.Value, &value, handlers.Cookies["current"], handlers.Cookies["previous"])
	// 	if err == nil {
	// 		data := handlers.FormatData{
	// 			Name:  value["name"],
	// 			Login: value["login"],
	// 			Info:  value["info"],
	// 		}
	// 		handlers.TempHome.Execute(w, data)
	// 	}
	// } else {
	// 	data := handlers.FormatData{
	// 		Message: "make a mistake and I will remember that о_о",
	// 	}
	// 	handlers.TempLogin.Execute(w, data)
	// }

}
