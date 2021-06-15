package data

import (
	"database/sql"
	"log"

	models "../Server/models"
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
	UserStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Utilisateur (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Pseudo 				STRING NOT NULL, 
		MAIL 				STRING UNIQUENOT NULL, 
		PASSWORD 			STRING NOT NULL, 
		InscriptionDate		STRING NOT NULL,
		NumberOfPosts 		INTEGER, 
		NumberOfFriend 		INTEGER,
		LastActivityDate 	STRING,
		Description 		STRING, 
		Image 				STRING
	  );`)
	UserStatement.Exec()
	CheckError(err, " statement line 33 ")
	UserStatement, err = Database.Prepare(`
	INSERT INTO Utilisateur (
		Pseudo, MAIL, PASSWORD, InscriptionDate, NumberOfPosts, NumberOfFriend, LastActivityDate, Description, Image
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	CheckError(err, " statement line 39 ")
	return UserStatement
}

func InitPostUserData() *sql.Stmt {
	PostStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS Posts (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT, 
		Titre 				STRING NOT NULL, 
		Description 		STRING NOT NULL, 
		Likes 				INTEGER, 
		Dislikes 			INTEGER, 
		Dates 				STRING NOT NULL, 
		AutorName 			STRING NOT NULL
	  );`)
	CheckError(err, "statement line 54")
	PostStatement.Exec()
	PostStatement, err = Database.Prepare(`
	INSERT INTO Posts (
	Titre, Description, Likes, Dislikes, Dates, AutorName
	) 
	VALUES (?, ?, ?, ?, ?, ?)`)
	CheckError(err, "statement line 61")
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
		Dislikes 		INTEGER, 
		Dates 			STRING
	);`)
	CheckError(err, "statement line 76")
	CommentStatement.Exec()
	CommentStatement, err = Database.Prepare(`
	INSERT INTO Commentaires (
		Comment, CommentPostID, User, Likes, Dislikes, Dates
	) 
	VALUES (?, ?, ?, ?, ?, ?)`)
	CheckError(err, "statement line 83")
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

func InitAllTopicUserData() *sql.Stmt {
	AllTopicStatement, err := Database.Prepare(`
	CREATE TABLE IF NOT EXISTS AllTopic (
		ID INTEGER PRIMARY KEY ASC AUTOINCREMENT,
		Name STRING NOT NULL, 
		NumberOfPublication INTEGER,
		LastActivity STRING NOT NULL
	)
	`)
	CheckError(err, "line 156 ")
	AllTopicStatement.Exec()

	AllTopicStatement, err = Database.Prepare(`
	INSERT INTO AllTopic (
		Name, NumberofPublication, LastActivity
	)
	VALUES (?, ?, ?)
	`)

	CheckError(err, "line 166 ")
	return AllTopicStatement

}

func InitAllTopicInDatabase() []models.Tag {
	var AllTopicStatement *sql.Stmt = InitAllTopicUserData()
	AllTags := []string{"Electronique", "Programmation", "Jeux", "Animation", "Musique", "Sport", "Ynov", "Horreur", "Film/Séries"}

	for _, Tag := range AllTags {
		AllTopicStatement.Exec(Tag, 0, "Never")
	}

	return ReturnAllTopic()
}

func MaybeHaveAlreadyTopicInDatabase() bool {
	rows, err := Database.Query("SELECT * FROM AllTopic")
	CheckError(err, "line 179 ")
	defer rows.Close()
	var boolean bool
	for rows.Next() {
		boolean = true
	}

	return boolean
}

func ReturnAllTopic() []models.Tag {
	if MaybeHaveAlreadyTopicInDatabase() {
		var AllTags []models.Tag
		rows, err := Database.Query("SELECT * FROM AllTopic")
		CheckError(err, "line 179 ")
		defer rows.Close()
		for rows.Next() {
			var CurrTag models.Tag
			var id int
			var Topic string
			var NumberOfPosts int
			var LastPostPublished string

			err = rows.Scan(&id, &Topic, &NumberOfPosts, &LastPostPublished)
			CheckError(err, "scan 197 ")
			CurrTag.Name = Topic
			CurrTag.NumberOfPosts = NumberOfPosts
			CurrTag.LastPostPublished = LastPostPublished
			AllTags = append(AllTags, CurrTag)

		}

		return AllTags
	} else {
		return InitAllTopicInDatabase()
	}

}
