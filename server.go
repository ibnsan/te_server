package main

import (
	"html/template"
	"net/http"
)

type formatData struct { //задаю формат для переменной которая у меня будет в страничке
	Text string
}

func main() {
	http.HandleFunc("/handler", handler) //собсно мой "обработчик"
	http.ListenAndServe(":80", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := formatData{ //указываю что нужно подгрузить в переменную на страничке
		Text: r.FormValue("mood"),
	}
	temp, _ := template.ParseFiles("page.html") //пасрю страницу
	temp.Execute(w, data)                       //передю страничке данные и вывожу ее
}

// test
