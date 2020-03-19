package handlers

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "auth",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, c)

	http.ServeFile(w, r, "/")
}
