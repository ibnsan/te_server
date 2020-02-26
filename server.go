package main

import (
	"html/template"
	"net/http"
)

type formatData struct {
	login    string
	password string
}

func main() {
	http.HandleFunc("/", check)
	http.HandleFunc("/handlerLogin", handlerLogin)
	http.HandleFunc("/handlerRegistration", handlerRegistration)
	http.ListenAndServe(":80", nil)
}

func check(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("pages/login.html") //если пользователь не залогинен - пересылаю его на страницу логина

	if len(r.Header["Cookie"]) != 0 && r.Header["Cookie"][0] == "auth=your_MD5_cookies" {
		temp, _ := template.ParseFiles("pages/home.html") //если пользователь залогинен - пересылаю его на главную страницу
	}
	temp.Execute(w, temp) //если пользователь залогинен, передам сюда его данные (к примеру логин в переменную), чтобы затем высветить
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	//здесь будет обработка данных при логине, проверка есть ли пользователь с таким логином и паролем
	data := formatData{
		login:    r.FormValue("login"),
		password: r.FormValue("pass"),
	}
	temp, _ := template.ParseFiles("page.html")
	temp.Execute(w, data)
}

func handlerRegistration(w http.ResponseWriter, r *http.Request) {
	//здесь обработка регистрации, проверка есть ли данные с таким логином\паролем, если такого логина нет тогда записываю его в бд и после направляю на главную страницу
	data := formatData{
		login: r.FormValue("mood"),
	}
	temp, _ := template.ParseFiles("page.html")
	temp.Execute(w, data)
}
