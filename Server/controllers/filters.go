package controllers

import (
	"fmt"
	"time"

	data "../../Database"
	models "../models"
)

func FilterPostsByTopic(Topic string, Posts []models.Post) []models.Post {
	if Topic == "AllTags" || Topic == "" {
		return Posts
	}
	var PostsFiltered []models.Post
	if Topic == "" {
		return Posts
	}
	for _, Post := range Posts {
		for _, CurrTopic := range Post.Tags {
			if CurrTopic.Tag == Topic {
				PostsFiltered = append(PostsFiltered, Post)
			}
		}
	}
	return PostsFiltered
}

// Tri par séléction
// ! Améliorer l'algorithme de tri (selection ==> Fusion)
func SortPostsByLike(form string, Posts []models.Post) []models.Post {
	if form == "" {
		return Posts
	}

	for i := 0; i < len(Posts); i++ {
		for j := i; j < len(Posts); j++ {
			if Posts[i].NbrLikes < Posts[j].NbrLikes {
				Posts[i], Posts[j] = Posts[j], Posts[i]
			}
		}
	}
	return Posts
}

func SortPostByDate(form string, Posts []models.Post) []models.Post {
	if form == "" {
		return Posts
	}

	for i := 0; i < len(Posts); i++ {
		for j := i; j < len(Posts); j++ {
			time_I, _ := time.Parse("02-01-2006 15:04:05", Posts[i].Dates)
			time_J, _ := time.Parse("02-01-2006 15:04:05", Posts[j].Dates)
			if time_I.Before(time_J) {
				Posts[i], Posts[j] = Posts[j], Posts[i]
			}
		}
	}

	return Posts
}

func SelectOneListOfComment(PostID string) []models.Comment {

	var UniqueListOfComment []models.Comment

	rows, err := data.Database.Query("SELECT * FROM Commentaires WHERE CommentPostID = " + PostID)
	CheckError(err, "rows line 71")
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
		CheckError(err, "scan line 94")

		var Commentaire models.Comment
		Commentaire.ID = ID
		Commentaire.PostID = CommentPostID
		Commentaire.Message = Comment
		Commentaire.NbrLikes = Likes
		Commentaire.NbrDislike = Dislikes
		Commentaire.Autor = User
		Commentaire.Dates = Dates

		UniqueListOfComment = append(UniqueListOfComment, Commentaire)
	}

	return UniqueListOfComment
}

func SelectOnePost(IDs string) []models.Post {
	if IDs == "" {
		return []models.Post{}
	}

	// fmt.Println(IDs)
	var UniquePost []models.Post
	rows, err := data.Database.Query("SELECT * from Posts WHERE ID = " + IDs)
	// fmt.Println(err, rows)
	CheckError(err, "rows line 65")
	defer rows.Close()

	for rows.Next() {
		// fmt.Println("HELLO")
		var ID int
		var title string
		var Description string
		var Likes int
		var Dislikes int
		var Dates string
		var AutorName string
		var CurrPost models.Post
		var CurrComment []models.Comment = SelectOneListOfComment(IDs)

		err = rows.Scan(&ID, &title, &Description, &Likes, &Dislikes, &Dates, &AutorName)
		CheckError(err, "scan 85")

		CurrPost.ID = ID
		CurrPost.Title = title
		CurrPost.Description = Description
		CurrPost.NbrLikes = Likes
		CurrPost.NbrDislike = Dislikes
		CurrPost.Autor = AutorName
		CurrPost.Dates = Dates
		CurrPost.Commentaires = CurrComment

		UniquePost = append(UniquePost, CurrPost)
	}
	fmt.Println(UniquePost)
	return UniquePost
}
