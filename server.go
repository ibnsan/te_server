package main

import (
	"database/sql"
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
	login   string
	pass    string
	Message string
}

func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
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

	http.HandleFunc("/", check)
	http.HandleFunc("/handlerLogin", handlerLogin)
	http.HandleFunc("/handlerRegistration", handlerRegistration)
	http.ListenAndServe(":80", nil)
	http.Handle("/js/", http.FileServer(http.Dir("path/to")))
}

func check(w http.ResponseWriter, r *http.Request) {

	if cookie, err := r.Cookie("auth"); err == nil {
		value := make(map[string]string)
		err = securecookie.DecodeMulti("auth", cookie.Value, &value, cookies["current"], cookies["previous"])
		if err == nil {
			// fmt.Fprintf(w, "The value of foo is %q", value["login"])
			temp, _ := template.ParseFiles("pages/home.html") //если пользователь залогинен - пересылаю его на главную страницу
			temp.Execute(w, temp)
		} else {
			data := formatData{
				Message: "make a mistake and I will remember that о_о",
			}
			temp, _ := template.ParseFiles("pages/login.html") //если пользователь не залогинен - пересылаю его на страницу логина
			temp.Execute(w, data)                              //если пользователь залогинен, передам сюда его данные (к примеру логин в переменную), чтобы затем высветить
		}
	}

	// if len(r.Header["Cookie"]) != 0 && r.Header["Cookie"][0] == "auth=your_MD5_cookies" {
	// 	temp, _ := template.ParseFiles("pages/home.html") //если пользователь залогинен - пересылаю его на главную страницу
	// 	temp.Execute(w, temp)
	// } else {
	// 	data := formatData{
	// 		Message: "make a mistake and I will remember that о_о",
	// 	}
	// 	temp, _ := template.ParseFiles("pages/login.html") //если пользователь не залогинен - пересылаю его на страницу логина
	// 	temp.Execute(w, data)                              //если пользователь залогинен, передам сюда его данные (к примеру логин в переменную), чтобы затем высветить
	// }

}

func handlerLogin(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "auser:12345678@/tes_bd")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	login := r.FormValue("login")
	password := r.FormValue("pass")

	row := db.QueryRow("select id from tes_bd.users where login = ? and pass = ?", login, password)

	p := formatData{}
	err = row.Scan(&p.id)

	if err != nil { //если логин и пароль ошибочны
		data := formatData{
			Message: "Oh, you were mistaken in the password or login ... oh my God what to do now =(",
		}
		temp, _ := template.ParseFiles("pages/login.html")
		temp.Execute(w, data)

	} else { //если логин и пароль верны
		// addCookie(w, "TestCookieName", "TestValue")

		value := map[string]string{
			"login": login,
		}
		if encoded, err := securecookie.EncodeMulti("auth", value, cookies["current"]); err == nil {
			cookie := &http.Cookie{
				Name:  "auth",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}

		temp, _ := template.ParseFiles("pages/home.html")
		temp.Execute(w, temp)
	}

	// test := p.id

	// if test != nil {
	// 	fmt.Println(p.id)
	// } else {
	// 	fmt.Println("Oh fuck you mean")
	// }

	// data := formatData{
	// 	login:    r.FormValue("login"),
	// 	pass: r.FormValue("pass"),
	// }

	// temp, _ := template.ParseFiles("page.html")
	// temp.Execute(w, data)
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
	//здесь обработка регистрации, проверка есть ли данные с таким логином\паролем, если такого логина нет тогда записываю его в бд и после направляю на главную страницу
	data := formatData{
		login: r.FormValue("mood"),
	}
	temp, _ := template.ParseFiles("page.html")
	temp.Execute(w, data)
}
