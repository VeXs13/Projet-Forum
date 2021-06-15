package pages

import (
	"database/sql"
	"fmt"
	"time"

	data "../../../Database"
	models "../../models"
)

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
		var Dislikes int
		var Dates string
		var AutorName string
		var CurrPost models.Post
		var CurrComment []models.Comment

		err = rows.Scan(&ID, &title, &Description, &Likes, &Dislikes, &Dates, &AutorName)
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
		var Dislikes int
		var Dates string

		err = rows.Scan(&ID, &Comment, &CommentPostID, &User, &Likes, &Dislikes, &Dates)
		CheckError(err, "scan line 944")

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
	WHERE Autor = '` + CurrUser.Pseudo + `' 
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
		statement, err := data.Database.Prepare("DELETE FROM " + Table + " WHERE Autor = '" + CurrUser.Pseudo + "' AND " + column + " = " + fmt.Sprint(ID))
		CheckError(err, "statement line 181")
		statement.Exec()
		return boolean
	}

	if Table == "Likes" {
		if Type == "Commentaires" {
			LikeStatement.Exec(CurrUser.Pseudo, 0, ID)
		} else {
			LikeStatement.Exec(CurrUser.Pseudo, ID, 0)
		}
	} else {
		if Type == "Commentaires" {
			DislikeStatement.Exec(CurrUser.Pseudo, 0, ID)
		} else {
			DislikeStatement.Exec(CurrUser.Pseudo, ID, 0)
		}
	}
	return boolean
}

func NewToogleLikeOrDislike(ID int, Type string, Table string) {
	// Table contains the table to modify in string (table Like, or Dislike)
	// Type Correspont either in Commentary, either at Posts (to like or dislike)
	if ID == 0 {
		return
	}

	// fmt.Println(ID)

	var statement *sql.Stmt
	// Select like from Post for instance
	rows, err := data.Database.Query("SELECT " + Table + " FROM " + Type + " WHERE ID = " + fmt.Sprint(ID))
	CheckError(err, "NEW ROWS line 193")
	defer rows.Close()
	var AncienLike int
	for rows.Next() {
		rows.Scan(&AncienLike)
	}

	if NewHaveLikedOrDisliked(ID, Type, Table) { // If the users Have already "toogle" the post
		// fmt.Println("true")
		// Init Statement
		statement, err = data.Database.Prepare("UPDATE " + Type + " SET " + Table + " = '" + fmt.Sprint(AncienLike-1) + "' WHERE ID = " + fmt.Sprint(ID))
	} else {
		// fmt.Println("false")
		statement, err = data.Database.Prepare("UPDATE " + Type + " SET " + Table + " = '" + fmt.Sprint(AncienLike+1) + "' WHERE ID = " + fmt.Sprint(ID))
	}

	CheckError(err, "New statement")
	statement.Exec()
	UpdateLastActivity() // Update Last activity of the Current User
}

func DeletePost(ID string) {
	if ID == "" {
		return
	}
	// Remove everything concerning the Post(likes, dislikes, comment, posts...)
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
	statement, err = data.Database.Prepare("UPDATE Utilisateur SET NumberOfPosts = NumberOfPosts - 1 WHERE ID = " + fmt.Sprint(CurrUser.ID))
	CheckError(err, "line code error 234")
	statement.Exec()
	CurrUser.NumberOfPosts--
	UpdateLastActivity() // Also update activity
}

func AddPost(PostsTitle, PostDescription string, ListofTags []string) {
	if PostsTitle == "" || PostDescription == "" {
		return
	}
	dt := time.Now()
	res, err := PostStatement.Exec(PostsTitle, Byte13To10(PostDescription), 0, 0, dt.Format("02-01-2006 15:04:05"), CurrUser.Pseudo)
	CheckError(err, "RESPONSE ")
	for _, Tag := range ListofTags {
		if Tag != "" {
			id, err := res.LastInsertId()
			CheckError(err, "RESPONSE ")
			TopicStatement.Exec(Tag, id)
			AddPostForTagsList(Tag)
		}
	}
	statement, err := data.Database.Prepare("UPDATE Utilisateur SET NumberOfPosts = NumberOfPosts + 1 WHERE ID = " + fmt.Sprint(CurrUser.ID))
	CurrUser.NumberOfPosts++
	CheckError(err, "line code 23444")
	statement.Exec()
	UpdateLastActivity() // Update activity
}

func AddPostForTagsList(Name string) {
	// Function will add in the column "Nulber of Post" of the AllTopic Table 1
	rows, err := data.Database.Query("SELECT NumberOfPublication FROM AllTopic WHERE NAME = '" + Name + "'")
	CheckError(err, "rows line 243")
	defer rows.Close()
	var AncienLike int
	for rows.Next() {
		rows.Scan(&AncienLike)
	}

	statement, err := data.Database.Prepare("UPDATE AllTopic SET NumberOfPublication = " + fmt.Sprint(AncienLike+1) + ", LastActivity = '" + time.Now().Format("02-01-2006 15:04:05") + "' WHERE Name = '" + Name + "'")
	CheckError(err, "line 250 rows ")
	statement.Exec()
}

func AddCommment(Comment string, PostID int) {
	if Comment == "" {
		return
	}
	dt := time.Now()
	CommentStatement.Exec(Byte13To10(Comment), PostID, CurrUser.Pseudo, 0, 0, dt.Format("02-01-2006 15:04:05"))
	UpdateLastActivity()
}

func UpdateLastActivity() {
	// Update LastActivity in DB
	statement, err := data.Database.Prepare("UPDATE Utilisateur SET LastActivityDate = '" + time.Now().Local().Format("02-01-2006 15:04:05") + "' WHERE ID = " + fmt.Sprint(CurrUser.ID))
	CheckError(err, "line 261 ")
	statement.Exec()
	CurrUser.LastActivityDate = time.Now().Local().Format("02-01-2006 15:04:05")
}

func DeleteComment(IDcomment string) {
	if IDcomment == "" {
		return
	}
	statement, err := data.Database.Prepare("DELETE FROM Commentaires WHERE ID = " + IDcomment)
	CheckError(err, "statement")
	statement.Exec()
	UpdateLastActivity()
}

func UpdateSession(guestValue string) {
	// Func will check if the user
	if guestValue != "" {
		CurrUser = models.User{}
	}
}
