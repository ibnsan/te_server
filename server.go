package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	http.Handle("/pages/style/", http.StripPrefix("/pages/style/", http.FileServer(http.Dir("pages/style"))))
	http.HandleFunc("/", checkAuth)
	http.HandleFunc("/handlerLogin", handlerLogin)
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/handlerRegistration", handlerRegistration)
	http.HandleFunc("/handlerLogout", handlerLogout)
	err := http.ListenAndServe(":8666", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) {

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

func handlerLogin(w http.ResponseWriter, r *http.Request) {
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

func registration(w http.ResponseWriter, r *http.Request) {
	data := formatData{
		Message: "",
	}
	tempRegister.Execute(w, data)
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
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

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "auth",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, c)

	http.ServeFile(w, r, "/")

	//тут есть проблема, почему то удалить куки можно только один раз, дальше они не удаляются, пока не понимаю почему (
}
