package main

import (
	"database/sql"
	"fmt"
	"html/template"
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

	http.HandleFunc("/", check)
	http.HandleFunc("/handlerLogin", handlerLogin)
	http.HandleFunc("/registration", registration)
	http.HandleFunc("/handlerRegistration", handlerRegistration)
	http.HandleFunc("/handlerLogout", handlerLogout)
	http.ListenAndServe(":80", nil)
}

func check(w http.ResponseWriter, r *http.Request) { //проверка авторизации

	if cookie, err := r.Cookie("auth"); err == nil { //если есть куки о том что пользователь залогинен - перекидываю его сразу на домашнюю страницу
		value := make(map[string]string)
		err = securecookie.DecodeMulti("auth", cookie.Value, &value, cookies["current"], cookies["previous"])
		if err == nil {
			data := formatData{
				Name:  value["name"],
				Login: value["login"],
				Info:  value["info"],
			}
			temp, _ := template.ParseFiles("pages/home.html")
			temp.Execute(w, data)

		}
	} else { //если нет куков о том что пользователь залогинен
		data := formatData{
			Message: "make a mistake and I will remember that о_о",
		}
		temp, _ := template.ParseFiles("pages/login.html")
		temp.Execute(w, data)
	}

}

func handlerLogin(w http.ResponseWriter, r *http.Request) { //обработка авторизации

	db, err := sql.Open("mysql", "auser:12345678@/tes_bd")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	login := r.FormValue("login")
	password := r.FormValue("pass")

	row := db.QueryRow("select id, name, info from tes_bd.users where login = ? and pass = ?", login, password)

	p := formatData{}
	err = row.Scan(&p.id, &p.Name, &p.Info)

	if err != nil { //если логин и пароль ошибочны
		data := formatData{
			Message: "Oh, you were mistaken in the password or login ... oh my God what to do now =(",
		}
		temp, _ := template.ParseFiles("pages/login.html")
		temp.Execute(w, data)

	} else { //если логин и пароль верны
		value := map[string]string{
			"login": login,
			"name":  p.Name,
			"info":  p.Info,
		}
		if encoded, err := securecookie.EncodeMulti("auth", value, cookies["current"]); err == nil {
			cookie := &http.Cookie{
				Name:  "auth",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie) //записываю данные в куки - можно конечно записать только логин или id, а потом по ним запросить все остальные данные
		}
		http.ServeFile(w, r, "/")
	}
}

func registration(w http.ResponseWriter, r *http.Request) { //загрузка страницы регистрации
	data := formatData{
		Message: "",
	}
	temp, _ := template.ParseFiles("pages/registration.html")
	temp.Execute(w, data)
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) { //обработка регистрации

	//получаю данные из формы
	login := r.FormValue("login")
	password := r.FormValue("pass")
	name := r.FormValue("name")
	info := r.FormValue("info")

	db, err := sql.Open("mysql", "auser:12345678@/tes_bd")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow("select id from tes_bd.users where login = ?", login)

	p := formatData{}
	err = row.Scan(&p.id)
	if err != nil { //проверяю, если пользователя с таким логином нет, регистрирую, если есть - пишу что нужен другой логин
		result, err := db.Exec("insert into tes_bd.users (login, pass, name, info) values (?, ?, ?, ?)", login, password, name, info)
		if err != nil {
			panic(err)
		} else {
			data := formatData{
				Message: "Yuhu, now you can log in (I hope you remember your password ...)",
			}
			temp, _ := template.ParseFiles("pages/login.html")
			temp.Execute(w, data) //если регистрация успешна, направляю на страницу логина
		}
		fmt.Println(result.LastInsertId())
		fmt.Println(result.RowsAffected())
	} else {
		data := formatData{
			Message: "Wow ... sorry bro, but this login is already busy with someone, please be smarter",
		}
		temp, _ := template.ParseFiles("pages/registration.html")
		temp.Execute(w, data)
	}

}

func handlerLogout(w http.ResponseWriter, r *http.Request) { //обработка выхода
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
