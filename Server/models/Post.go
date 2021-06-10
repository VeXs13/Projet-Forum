package models

type Comment struct {
	ID         int
	PostID     int
	Message    string
	NbrLikes   int
	NbrDislike int
	Autor      string
	Dates      string
}

type LikeOrDislike struct {
	UserID    int
	PostID    int
	CommentID int
}

type Post struct {
	ID           int
	Title        string
	Description  string
	NbrLikes     int
	NbrDislike   int
	Commentaires []Comment
	Autor        string
	Tags         []Tags
	Dates        string
}

type Tags struct {
	Tag    string
	PostID int
}
