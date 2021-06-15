package pages

import (
	"html/template"
	"net/http"

	controllers ".."
	models "../../models"
)

func PostsPage(w http.ResponseWriter, r *http.Request) {
	if CurrUser.Session == "Connected" {
		PostID := r.FormValue("DeletePost")
		PostTitle := r.FormValue("PostTitle")
		PostDescription := r.FormValue("PostDescription")
		ListofTags := []string{r.FormValue("Jeux"), r.FormValue("Animation"), r.FormValue("Electronique"), r.FormValue("Programmation"), r.FormValue("Musique"), r.FormValue("Sport"), r.FormValue("Ynov"), r.FormValue("Horreur"), r.FormValue("Film/Séries")}
		AddPost(PostTitle, PostDescription, ListofTags)
		DeletePost(PostID)
	}
	ActualizePosts()
	var UsersFiltered models.User = CurrUser
	TopicSelectedInHomePage := r.FormValue("TagsSelection")
	FilterByDates := r.FormValue("FilterByDates")
	FilterByLikes := r.FormValue("FilterByLikes")

	UsersFiltered.Posts = controllers.FilterPostsByTopic(TopicSelectedInHomePage, UsersFiltered.Posts)
	UsersFiltered.Posts = controllers.SortPostsByLike(FilterByLikes, UsersFiltered.Posts)
	UsersFiltered.Posts = controllers.SortPostByDate(FilterByDates, UsersFiltered.Posts)

	// ==> Récuperer les fonctions de filtres
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./Client/PostLists.html"))
	err := t.ExecuteTemplate(w, "PostLists.html", UsersFiltered)
	CheckError(err, "template Post")
}
