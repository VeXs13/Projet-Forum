package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf(message+" execution : %s", err)
	}
}

var Database *sql.DB

func InitUserDatabase() *sql.Stmt {
	Database, _ = sql.Open("sqlite3", "./Database/UserT.db")
	UserStatement, _ := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Utilisateur (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Nom 		STRING NOT NULL, 
		PRENOM 		STRING NOT NULL, 
		MAIL 		STRING UNIQUENOT NULL, 
		PASSWORD 	STRING NOT NULL, 
		ADDRESSE	STRING NOT NULL
	  );`)
	UserStatement.Exec()
	UserStatement, _ = Database.Prepare(`
	INSERT INTO Utilisateur (
	Nom, PRENOM, MAIL, PASSWORD, ADDRESSE
	) 
	VALUES (?, ?, ?, ?, ?)`)
	return UserStatement
}

func InitPostUserData() *sql.Stmt {
	PostStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Posts (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Titre 				STRING NOT NULL, 
		Description 		STRING NOT NULL, 
		Likes 				INTEGER, 
		WhoHaveLiked 		STRING NOT NULL, 
		Dislikes 			INTEGER, 
		WhoHaveDisliked 	STRING, 
		Dates 				STRING NOT NULL, 
		AutorName 			STRING NOT NULL
	  );`)
	CheckError(err, "statement")
	PostStatement.Exec()
	PostStatement, err = Database.Prepare(`
	INSERT INTO Posts (
	Titre, Description, Likes, WhoHaveLiked, Dislikes, WhoHaveDisliked, Dates, AutorName
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	CheckError(err, "statement")
	return PostStatement
}

func InitCommentData() *sql.Stmt {
	CommentStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Commentaires (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Comment 		STRING NOT NULL,
		CommentPostID 	INTEGER, 
		User 			STRING NOT NULL, 
		Likes 			INTEGER, 
		WhoHaveLiked 	STRING, 
		Dislikes 		INTEGER, 
		WhoHaveDisliked STRING, 
		Dates 			STRING
	);`)
	CheckError(err, "statement")
	CommentStatement.Exec()
	CommentStatement, err = Database.Prepare(`
	INSERT INTO Commentaires (
		Comment, CommentPostID, User, Likes, WhoHaveLiked, Dislikes, WhoHaveDisliked, Dates
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	CheckError(err, "statement")
	return CommentStatement
}

func InitTopicData() *sql.Stmt {
	TopicStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Topics (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Topic 		STRING NOT NULL,
		PostID 		INTEGER 
	);`)
	CheckError(err, "Topic Statement")
	TopicStatement.Exec()
	TopicStatement, err = Database.Prepare(`
	INSERT INTO Topics (
		Topic, PostID
	)
	VALUES (?, ?)`)
	CheckError(err, "Topic Statement")
	return TopicStatement
}

func InitDislikeData() *sql.Stmt {
	DislikeStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Dislikes (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT,
		Autor STRING NOT NULL,
		PostID INTEGER,
		CommentID INTEGER 
	)
	`)

	CheckError(err, "DislikeStatement")
	DislikeStatement.Exec()

	DislikeStatement, err = Database.Prepare(`
	INSERT INTO Dislikes (
		Autor, PostID, CommentID
	)
	VALUES (?, ?, ?)`)
	CheckError(err, "FavPostStatement")
	return DislikeStatement
}

func InitLikeData() *sql.Stmt {
	LikeStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Likes (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT,
		Autor STRING NOT NULL,
		PostID INTEGER,
		CommentID INTEGER 
	)
	`)

	CheckError(err, "LikeStatement")
	LikeStatement.Exec()

	LikeStatement, err = Database.Prepare(`
	INSERT INTO Likes (
		Autor, PostID, CommentID
	)
	VALUES (?, ?, ?)`)
	CheckError(err, "FavPostStatement")
	return LikeStatement
}
