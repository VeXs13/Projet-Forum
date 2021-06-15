package pages

import (
	"html/template"
	"net/http"
)

func FriendPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./Client/Friend.html"))
	err := t.ExecuteTemplate(w, "Friend.html", CurrUser)
	CheckError(err, "template")
}
