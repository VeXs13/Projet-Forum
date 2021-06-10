package pages

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http" // Needed to use Split
	"time"

	controllers ".."
	data "../../../Database"
	models "../../models"
	_ "github.com/mattn/go-sqlite3"
)

var UserStatement *sql.Stmt = data.InitUserDatabase()
var PostStatement *sql.Stmt = data.InitPostUserData()
var CommentStatement *sql.Stmt = data.InitCommentData()
var TopicStatement *sql.Stmt = data.InitTopicData()
var LikeStatement *sql.Stmt = data.InitLikeData()
var DislikeStatement *sql.Stmt = data.InitDislikeData()

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf(message+" execution : %s", err)
	}
}

func Byte13To10(text string) string {
	var Newtext string
	for _, char := range text {
		if byte(char) == 13 {
			Newtext += "\n"
		} else {
			Newtext += string(char)
		}
	}

	return Newtext
}

func ActualizeTags() []models.Tags {
	var Topics []models.Tags

	rows, err := data.Database.Query("SELECT * FROM Topics")
	CheckError(err, "database query ")
	defer rows.Close()

	for rows.Next() {
		var ID int
		var Topic string
		var PostID int
		var CurrTopic models.Tags

		err = rows.Scan(&ID, &Topic, &PostID)
		CheckError(err, "scan")
		CurrTopic.PostID = PostID
		CurrTopic.Tag = Topic
		Topics = append(Topics, CurrTopic)
	}
	// fmt.Println(Topics)
	return Topics
}

func ActualizePosts() {
	CommentSlice := ActualizeComment()
	TagsSlice := ActualizeTags()
	// fmt.Println(TagsSlice)
	var Posts []models.Post
	rows, err := data.Database.Query("SELECT * FROM Posts")
	CheckError(err, "database")
	defer rows.Close()

	for rows.Next() {
		var ID int
		var title string
		var Description string
		var Likes int
		var UserLiked string
		var Dislikes int
		var WhoHaveDisliked string
		var Dates string
		var AutorName string
		var CurrPost models.Post
		var CurrComment []models.Comment

		err = rows.Scan(&ID, &title, &Description, &Likes, &UserLiked, &Dislikes, &WhoHaveDisliked, &Dates, &AutorName)
		CheckError(err, "scan")

		CurrPost.ID = ID
		CurrPost.Title = title
		CurrPost.Description = Description
		CurrPost.NbrLikes = Likes
		CurrPost.NbrDislike = Dislikes
		CurrPost.Autor = AutorName
		CurrPost.Dates = Dates
		//  Optimiser la boucle en supprimant les élements ajouté
		for _, Comment := range CommentSlice {
			if Comment.PostID == CurrPost.ID {
				CurrComment = append(CurrComment, Comment)
			}
		}

		for _, Tag := range TagsSlice {
			if Tag.PostID == ID {
				CurrPost.Tags = append(CurrPost.Tags, Tag)
			}
		}
		CurrPost.Commentaires = append(CurrPost.Commentaires, CurrComment...)
		Posts = append(Posts, CurrPost)
	}
	CurrUser.Posts = Posts
}

func ActualizeComment() []models.Comment {
	var CommentSlice []models.Comment
	rows, err := data.Database.Query("SELECT * FROM Commentaires")
	CheckError(err, "database")
	defer rows.Close()

	for rows.Next() {
		var ID int
		var Comment string
		var CommentPostID int
		var User string
		var Likes int
		var WhoHaveLiked string
		var Dislikes int
		var WhoHaveDisliked string
		var Dates string

		err = rows.Scan(&ID, &Comment, &CommentPostID, &User, &Likes, &WhoHaveLiked, &Dislikes, &WhoHaveDisliked, &Dates)
		CheckError(err, "scan line 94")

		var Commentaire models.Comment
		Commentaire.ID = ID
		Commentaire.PostID = CommentPostID
		Commentaire.Message = Comment
		Commentaire.NbrLikes = Likes
		Commentaire.NbrDislike = Dislikes
		Commentaire.Autor = User
		Commentaire.Dates = Dates

		CommentSlice = append(CommentSlice, Commentaire)
	}
	return CommentSlice
}

func NewHaveLikedOrDisliked(ID int, Type string, Table string) bool {
	var boolean bool
	var column string
	if Type == "Commentaires" {
		column = "CommentID"
	} else {
		column = "PostID"
	}
	// On cherche dans notre base de données si l'utilisateur a déjà liké/disliké ou non le Post/Commentaire en question
	rows, err := data.Database.Query(`
	SELECT * from ` + Table + ` 
	WHERE Autor = '` + CurrUser.Name + `' 
	AND ` + column + ` = '` + fmt.Sprint(ID) + `'
	`)
	CheckError(err, "rows from New error")
	defer rows.Close()
	// Si le resultat de la réponse est itérable, c'est que l'utilisateur a déjà liké/disliké. On retourne true
	for rows.Next() {
		boolean = true
	}
	// Sinon on retourne false car il n'a pas été liké
	if boolean {
		statement, err := data.Database.Prepare("DELETE FROM " + Table + " WHERE Autor = '" + CurrUser.Name + "' AND " + column + " = " + fmt.Sprint(ID))
		CheckError(err, "statement line 181")
		statement.Exec()
		return boolean
	}

	if Table == "Likes" {
		if Type == "Commentaires" {
			LikeStatement.Exec(CurrUser.Name, 0, ID)
		} else {
			LikeStatement.Exec(CurrUser.Name, ID, 0)
		}
	} else {
		if Type == "Commentaires" {
			DislikeStatement.Exec(CurrUser.Name, 0, ID)
		} else {
			DislikeStatement.Exec(CurrUser.Name, ID, 0)
		}
	}
	return boolean
}

func NewToogleLikeOrDislike(ID int, Type string, Table string) {
	// Table contient le tableau à modifier en string (soit le tableau Like, soit le tableau Dislike)
	// Type Correspont soit au Commentaires, soit au Post (à liker ou disliker)
	if ID == 0 {
		return
	}

	// fmt.Println(ID)

	var statement *sql.Stmt
	rows, err := data.Database.Query("SELECT " + Table + " FROM " + Type + " WHERE ID = " + fmt.Sprint(ID))
	CheckError(err, "NEW ROWS line 193")
	defer rows.Close()
	var AncienLike int
	for rows.Next() {
		rows.Scan(&AncienLike)
	}

	if NewHaveLikedOrDisliked(ID, Type, Table) {
		fmt.Println("true")
		statement, err = data.Database.Prepare("UPDATE " + Type + " SET " + Table + " = '" + fmt.Sprint(AncienLike-1) + "' WHERE ID = " + fmt.Sprint(ID))
	} else {
		fmt.Println("false")
		statement, err = data.Database.Prepare("UPDATE " + Type + " SET " + Table + " = '" + fmt.Sprint(AncienLike+1) + "' WHERE ID = " + fmt.Sprint(ID))
	}

	CheckError(err, "New statement")
	statement.Exec()
}

func DeletePost(ID string) {
	if ID == "" {
		return
	}
	statement, err := data.Database.Prepare("DELETE FROM Posts WHERE ID = " + ID)
	CheckError(err, "statement")
	statement.Exec()
	statement, err = data.Database.Prepare("DELETE FROM Commentaires WHERE CommentPostID = " + ID)
	CheckError(err, "statement")
	statement.Exec()
	statement, err = data.Database.Prepare("DELETE FROM Likes WHERE PostID = " + ID)
	CheckError(err, "DELETE")
	statement.Exec()
	statement, err = data.Database.Prepare("DELETE FROM Dislikes WHERE PostID = " + ID)
	CheckError(err, "DELETE")
	statement.Exec()
}

func AddPost(PostsTitle, PostDescription string, ListofTags []string) {
	if PostsTitle == "" || PostDescription == "" {
		return
	}
	dt := time.Now()
	res, err := PostStatement.Exec(PostsTitle, Byte13To10(PostDescription), 0, "", 0, "", dt.Format("01-02-2006 15:04:05"), CurrUser.Name)
	CheckError(err, "RESPONSE ")
	for _, Tag := range ListofTags {
		if Tag != "" {
			id, err := res.LastInsertId()
			CheckError(err, "RESPONSE ")
			/* L'ID commençant par 1, L'ID de la publication crée correspont à la longueur de la liste de post, plus 1 */
			TopicStatement.Exec(Tag, id)
		}
	}
}

func AddCommment(Comment, Dates string, PostID int) {
	if Comment == "" {
		return
	}
	dt := time.Now()
	CommentStatement.Exec(Byte13To10(Comment), PostID, CurrUser.Name, 0, "", 0, "", dt.Format("01-02-2006 15:04:05"))
}

func DeleteComment(IDcomment string) {
	if IDcomment == "" {
		return
	}
	statement, err := data.Database.Prepare("DELETE FROM Commentaires WHERE ID = " + IDcomment)
	CheckError(err, "statement")
	statement.Exec()
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if len(CurrUser.Name) != 0 {
		DelPost := r.FormValue("DeletePost")
		LikedPost := r.FormValue("Like")
		DislikedPost := r.FormValue("DislikedPost")
		PostTitle := r.FormValue("PostTitle")
		PostDescription := r.FormValue("PostDescription")
		DisconnectButton := r.FormValue("disconnect")
		Comment := r.FormValue("Commentaire")
		PostID := r.FormValue("PostID")
		Date := r.FormValue("Date")
		DelComment := r.FormValue("DeleteComment")
		LikedComment := r.FormValue("LikeComment")
		DislikeComment := r.FormValue("DislikeComment")
		ListofTags := []string{r.FormValue("Informatique"), r.FormValue("Jeux"), r.FormValue("Animation"), r.FormValue("Golang")}
		DeleteComment(DelComment)
		AddCommment(Comment, Date, controllers.ConvertToInt(PostID))
		DeletePost(DelPost)
		AddPost(PostTitle, PostDescription, ListofTags)
		Disconnect(w, r, DisconnectButton)
		NewToogleLikeOrDislike(controllers.ConvertToInt(LikedComment), "Commentaires", "Likes")
		NewToogleLikeOrDislike(controllers.ConvertToInt(LikedPost), "Posts", "Likes")
		NewToogleLikeOrDislike(controllers.ConvertToInt(DislikedPost), "Posts", "Dislikes")
		NewToogleLikeOrDislike(controllers.ConvertToInt(DislikeComment), "Commentaires", "Dislikes")
		ActualizePosts()
		t := template.New("Welcome")
		t = template.Must(t.ParseFiles("./Client/Home.html"))
		err := t.ExecuteTemplate(w, "Home.html", CurrUser)

		CheckError(err, "template")
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}
