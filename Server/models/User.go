package models

type User struct {
	// password         string
	ID               int
	Pseudo           string
	Mail             string
	InscriptionDate  string
	NumberOfPosts    int
	NumberOfFriend   int
	LastActivityDate string
	Description      string
	Image            string
	Posts            []Post
	AllTags          []Tag
	Session          string
	ErrMessage       string
}

type Tag struct {
	Name              string
	NumberOfPosts     int
	LastPostPublished string
}
