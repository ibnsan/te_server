package main

import (
	"log"
	"net/http"

	"github.com/ibnsan/te_server/handlers"
)

func main() {

	http.Handle("/pages/style/", http.StripPrefix("/pages/style/", http.FileServer(http.Dir("pages/style"))))
	http.HandleFunc("/", handlers.Сheckauth)
	http.HandleFunc("/handlerLogin", handlers.HandlerLogin)
	http.HandleFunc("/registration", handlers.NewRegistration)
	http.HandleFunc("/handlerRegistration", handlers.Registration)
	http.HandleFunc("/handlerLogout", handlers.Logout)
	err := http.ListenAndServe(":8666", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
