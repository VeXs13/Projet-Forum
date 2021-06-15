package main

import (
	"fmt"
	"log"
	"net/http"

	pages "./Server/controllers/pages"

	_ "github.com/mattn/go-sqlite3"
	// "github.com/gorilla/pat"
	// "github.com/gorilla/sessions"
	// "github.com/markbates/goth"
	// "github.com/markbates/goth/gothic"
	// "github.com/markbates/goth/providers/google"
)

/* func someHandler(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:  "theme",
		Value: "dark",
	}
	http.SetCookie(w, &c)
} */

func main() {

	// key := "77LJcoQ1jogSqmEZSZr2mPJB" // Replace with your SESSION_SECRET or similar
	// maxAge := 86400 * 30              // 30 days
	// isProd := false                   // Set to true when serving over https

	// store := sessions.NewCookieStore([]byte(key))
	// store.MaxAge(maxAge)
	// store.Options.Path = "/"
	// store.Options.HttpOnly = true // HttpOnly should always be enabled
	// store.Options.Secure = isProd

	// gothic.Store = store

	// goth.UseProviders(
	// 	google.New("413087465486-ifufkb7hjv4ivn2uqqkecedto6gr2v9p.apps.googleusercontent.com", "77LJcoQ1jogSqmEZSZr2mPJB", "http://localhost:3001/auth/google/callback", "email", "profile"),
	// )

	// p := pat.New()
	// p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

	// 	user, err := gothic.CompleteUserAuth(res, req)
	// 	if err != nil {
	// 		fmt.Fprintln(res, err)
	// 		return
	// 	}
	// 	t, _ := template.ParseFiles("Client/index.html")
	// 	fmt.Println(user)
	// 	t.Execute(res, user)
	// })

	// p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
	// 	gothic.BeginAuthHandler(res, req)
	// })

	// p.Get("/", func(res http.ResponseWriter, req *http.Request) {
	// 	t, _ := template.ParseFiles("Client/Home.html")
	// 	t.Execute(res, false)
	// })
	fs := http.FileServer(http.Dir("./Client/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/login", pages.LoginPage)
	http.HandleFunc("/", pages.HomePage)
	http.HandleFunc("/profil", pages.ProfilPage)
	http.HandleFunc("/Posts", pages.PostsPage)
	http.HandleFunc("/Posts/Details", pages.PostDetailPage)
	http.HandleFunc("/friend", pages.FriendPage)
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
			 Session : Active/Expired/Disconnected
		 } Else {

		 }
*/
// ? Name : "User"
// ? Value : User.Name (encoded)
