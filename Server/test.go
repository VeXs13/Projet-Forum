package main

import (
	"fmt"
	"log"
	"net/http"

	pages "./controllers/pages"
	_ "github.com/mattn/go-sqlite3"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:  "theme",
		Value: "dark",
	}
	http.SetCookie(w, &c)
}

func main() {
	fs := http.FileServer(http.Dir("../Client/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/login", pages.LoginPage)
	http.HandleFunc("/", pages.HomePage)
	http.HandleFunc("/test", someHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
