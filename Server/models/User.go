package models

type User struct {
	ID        int
	Pseudo    string
	password  string
	Mail      string
	Name      string
	LastName  string
	Age       string
	Posts     []Post
	Connected bool
}
