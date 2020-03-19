package handlers

import "net/http"

func NewRegistration(w http.ResponseWriter, r *http.Request) {
	data := formatData{
		Message: "",
	}
	tempRegister.Execute(w, data)
}
