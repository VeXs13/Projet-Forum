package pages

import (
	"fmt"
	"html/template"
	"net/http"

	controllers ".."
)

func PostDetailPage(w http.ResponseWriter, r *http.Request) {
	DisconnectButton := r.FormValue("disconnect")
	Comment := r.FormValue("Commentaire")
	PostID := r.FormValue("PostID")
	DelComment := r.FormValue("DeleteComment")
	LikedPost := r.FormValue("Like")
	DislikedPost := r.FormValue("DislikedPost")
	LikedComment := r.FormValue("LikeComment")
	DislikeComment := r.FormValue("DislikeComment")
	Disconnect(w, r, DisconnectButton)
	DeleteComment(DelComment)
	AddCommment(Comment, controllers.ConvertToInt(PostID))
	NewToogleLikeOrDislike(controllers.ConvertToInt(LikedComment), "Commentaires", "Likes")
	NewToogleLikeOrDislike(controllers.ConvertToInt(LikedPost), "Posts", "Likes")
	NewToogleLikeOrDislike(controllers.ConvertToInt(DislikedPost), "Posts", "Dislikes")
	NewToogleLikeOrDislike(controllers.ConvertToInt(DislikeComment), "Commentaires", "Dislikes")
	ActualizePosts()
	PostIDSelected := r.FormValue("PostID")
	fmt.Println(PostIDSelected)
	User := CurrUser
	User.Posts = controllers.SelectOnePost(PostIDSelected)
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./Client/PostDetail.html"))
	err := t.ExecuteTemplate(w, "PostDetail.html", User)
	CheckError(err, "template Post")
}
