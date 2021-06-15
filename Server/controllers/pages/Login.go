package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	controllers ".."
	data "../../../Database"
	models "../../models"

	_ "github.com/mattn/go-sqlite3"
)

var CurrUser models.User = models.User{AllTags: data.ReturnAllTopic()}
var ErrMessage string // Err message which display the potential message error

func LoginCheck(email string, password string) bool {
	// Login check function which will check if the Email and password are correct
	if email == "" || password == "" { // If empty form
		return false
	}
	// Query to recup the user with the email
	rows, err := data.Database.Query("SELECT * FROM Utilisateur WHERE MAIL='" + email + "'")
	CheckError(err, "database")
	defer rows.Close()
	// if we enter
	for rows.Next() {
		var uid int
		var pseudo string
		var mail string
		var passwd string
		var InscriptionDate string
		var NumberOfPosts int
		var NumberOfFriend int
		var LastActivityDate string
		var Description string
		var Image string
		err = rows.Scan(&uid, &pseudo, &mail, &passwd, &InscriptionDate, &NumberOfPosts, &NumberOfFriend, &LastActivityDate, &Description, &Image)
		CheckError(err, "row")
		if controllers.CheckPasswordHash(password, passwd) {
			// Rediriger vers la page d'acceuil en étant connecté avec le compte
			CurrUser.ID = uid
			CurrUser.Mail = mail
			CurrUser.Pseudo = pseudo
			CurrUser.Session = "Connected"
			CurrUser.Description = Description
			CurrUser.Image = Image
			CurrUser.LastActivityDate = LastActivityDate
			CurrUser.InscriptionDate = InscriptionDate
			CurrUser.NumberOfPosts = NumberOfPosts
			CurrUser.NumberOfFriend = NumberOfFriend
			// fmt.Println(CurrUser)
			return true
		} else {
			// Dipslay "mot de passe ou email incorrect"
			ErrMessage = "Mot de passe ou Email inccorect"
			fmt.Println("ACCEES REFUSÉ, MDP INCORRECT")
			return false
		}
	}

	// Rester sur la page, en affichant "mot de passe ou email incorrect"
	fmt.Println("ADRESSE EMAIL INTROUVABLE")
	ErrMessage = "Mot de passe ou Email Incorrect"
	return false
}

func RegisterCheck(Name, Mail string) (bool, string) {
	rows, err := data.Database.Query("SELECT * FROM Utilisateur WHERE MAIL='" + Mail + "'")
	CheckError(err, "Rows")
	defer rows.Close()
	if rows.Next() {
		return true, "Cette Adresse Email est déjà utilisé"
	}

	rows, err = data.Database.Query("SELECT * FROM Utilisateur WHERE Pseudo='" + Name + "'")
	CheckError(err, "Rows")
	if rows.Next() {
		return true, "Ce nom est déjà utilisé"
	}

	return false, ""
}

func Register(Name, Mail, password string) {
	if Name == "" || password == "" || Mail == "" {
		return
	}
	IsExistNameOrEmail, errMsg := RegisterCheck(Name, Mail)
	if IsExistNameOrEmail {
		ErrMessage = errMsg
	} else {
		password, err := controllers.HashPassword(password)
		CheckError(err, "Hashing")
		UserStatement.Exec(Name, Mail, password, time.Now().Local().Format("02-01-2006"), 0, 0, "Never", "", "") // On ajoute à la base de donnée les informations
		ErrMessage = "Votre inscription a bien été prit en compte"
	}
}

func InitCookie(w http.ResponseWriter, r *http.Request) {
	CookieLastLogged := &http.Cookie{
		Name:  "LastLogged",
		Value: strconv.FormatInt(time.Now().Unix(), 10),
	}
	CookieUser := &http.Cookie{
		Name:  "UserID",
		Value: strconv.FormatInt(int64(CurrUser.ID), 10),
	}

	http.SetCookie(w, CookieLastLogged)
	http.SetCookie(w, CookieUser)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// If the User want to reconnect
	if r.FormValue("Reconnect") != "" {
		CurrUser = models.User{} // The previous account will be deleted
	}
	// Recup data of form
	emailConnexion := r.FormValue("emailConnexion")
	passwordConnexion := r.FormValue("passwordConnexion")
	Name := r.FormValue("Name")
	emailInscription := r.FormValue("emailInscription")
	passwordInscription := r.FormValue("passwordInscription")
	Register(Name, emailInscription, passwordInscription) // Register Check
	if LoginCheck(emailConnexion, passwordConnexion) {    // If login and passwd are valid
		ErrMessage = ""
		InitCookie(w, r)
		fmt.Println(CurrUser)
		fmt.Println(r.Cookie("LastLogged"))
		http.Redirect(w, r, "/", 301)
	}
	t := template.New("Authentification")
	t = template.Must(t.ParseFiles("./Client/test.html"))
	err := t.ExecuteTemplate(w, "test.html", ErrMessage)
	CheckError(err, "template")
}
