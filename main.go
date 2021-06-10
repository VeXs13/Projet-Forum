package main

import (
	"fmt"
	"log"
	"net/http"

	pages "./Server/controllers/pages"

	_ "github.com/mattn/go-sqlite3"
)

/* func someHandler(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:  "theme",
		Value: "dark",
	}
	http.SetCookie(w, &c)
} */

func main() {
	fs := http.FileServer(http.Dir("../Client/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/login", pages.LoginPage)
	http.HandleFunc("/", pages.HomePage)
	// http.HandleFunc("/test", someHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

/* package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("LastLogged")
	cookie := &http.Cookie{
		Name:  "LastLogged",
		Value: strconv.FormatInt(time.Now().Unix(), 10),
	}
	http.SetCookie(w, cookie)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Hello, you're here for the first time")
	} else {
		timeint, _ := strconv.ParseInt(c.Value, 10, 0)
		fmt.Fprintf(w, "Hi, your last visit was at "+time.Unix(timeint, 0).Format("15:04:05"))
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))
} */

// * Name : "LastLogged"
// * Value : Date
/*
	 if DiffDates > 5h : {
		 launch HTML {
			 State : Disconnected
			 Objective : Connect
		 } Else {

		 }
*/
// ? Name : "User"
// ? Value : User.Name (encoded)
