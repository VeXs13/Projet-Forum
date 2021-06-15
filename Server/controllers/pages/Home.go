package pages

import (
	"database/sql"
	"html/template"
	"log"
	"net/http" // Needed to use Split

	data "../../../Database"
	_ "github.com/mattn/go-sqlite3"
)

var UserStatement *sql.Stmt = data.InitUserDatabase()
var PostStatement *sql.Stmt = data.InitPostUserData()
var CommentStatement *sql.Stmt = data.InitCommentData()
var TopicStatement *sql.Stmt = data.InitTopicData()
var LikeStatement *sql.Stmt = data.InitLikeData()
var DislikeStatement *sql.Stmt = data.InitDislikeData()
var AllTopicStatement *sql.Stmt = data.InitAllTopicUserData()

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf(message+" execution : %s", err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if CurrUser.Session == "Connected" {
		// We Check the potentials action autorized by the connected user
		DelPost := r.FormValue("DeletePost")
		PostTitle := r.FormValue("PostTitle")
		PostDescription := r.FormValue("PostDescription")
		DisconnectButton := r.FormValue("disconnect")
		ListofTags := []string{r.FormValue("Jeux"), r.FormValue("Animation"), r.FormValue("Electronique"), r.FormValue("Programmation"), r.FormValue("Musique"), r.FormValue("Sport"), r.FormValue("Ynov"), r.FormValue("Horreur"), r.FormValue("Film/Séries")}
		Disconnect(w, r, DisconnectButton)
		DeletePost(DelPost)
		AddPost(PostTitle, PostDescription, ListofTags)
	}
	CurrUser.AllTags = data.ReturnAllTopic()
	ActualizePosts() // Actualize Posts from the Other Users which also published or commented something
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./Client/index.html"))
	err := t.ExecuteTemplate(w, "index.html", CurrUser)
	CheckError(err, "template")
}
