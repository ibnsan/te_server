package handlers

import (
	"fmt"
	"net/http"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("pass")
	name := r.FormValue("name")
	info := r.FormValue("info")

	row := db.QueryRow("select id from tes_bd.users where login = ?", login)

	p := formatData{}
	dataDB := row.Scan(&p.id)

	if dataDB != nil {
		result, err := db.Exec("insert into tes_bd.users (login, pass, name, info) values (?, ?, ?, ?)", login, password, name, info)
		if err != nil {
			panic(err)
		} else {
			data := formatData{
				Message: "Yuhu, now you can log in (I hope you remember your password ...)",
			}
			tempLogin.Execute(w, data)
		}
		fmt.Println(result.LastInsertId())
		fmt.Println(result.RowsAffected())
	} else {
		data := formatData{
			Message: "Wow ... sorry bro, but this login is already busy with someone, please be smarter",
		}
		tempRegister.Execute(w, data)
	}

}
