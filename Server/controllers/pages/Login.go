package pages

import (
	"fmt"
	"html/template"
	"net/http"

	controllers ".."
	data "../../../Database"
	models "../../models"

	_ "github.com/mattn/go-sqlite3"
)

var CurrUser models.User
var ErrMessage string

func LoginCheck(email string, password string) bool {
	if email == "" || password == "" {
		return false
	}
	rows, err := data.Database.Query("SELECT * FROM Utilisateur WHERE MAIL='" + email + "'")
	CheckError(err, "database")
	defer rows.Close()
	for rows.Next() {
		var nom string
		var prenom string
		var mail string
		var passwd string
		var adress string
		var uid string
		err = rows.Scan(&uid, &nom, &prenom, &mail, &passwd, &adress)
		CheckError(err, "row")
		if controllers.CheckPasswordHash(password, passwd) {
			// Rediriger vers la page d'acceuil en étant connecté avec le compte
			CurrUser.ID = controllers.ConvertToInt(uid)
			CurrUser.Mail = mail
			CurrUser.Name = prenom
			// fmt.Println(CurrUser)
			return true
		} else {
			// Rester sur la page, en affichant "mot de passe ou email incorrect"
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
	for rows.Next() {
		return true, "Cette Adresse Email est déjà utilisé"
	}

	rows, err = data.Database.Query("SELECT * FROM Utilisateur WHERE PRENOM='" + Name + "'")
	CheckError(err, "Rows")
	for rows.Next() {
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
		UserStatement.Exec("", Name, Mail, password, "") // On ajoute à la base de donnée les informations
		ErrMessage = "Votre inscription a bien été prit en compte"
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// Récuperer les données du formulaire
	emailConnexion := r.FormValue("emailConnexion")
	passwordConnexion := r.FormValue("passwordConnexion")
	Name := r.FormValue("Name")
	emailInscription := r.FormValue("emailInscription")
	passwordInscription := r.FormValue("passwordInscription")
	Register(Name, emailInscription, passwordInscription)
	if LoginCheck(emailConnexion, passwordConnexion) { // On vérifie que les mdp et email sont valides
		ErrMessage = ""
		http.Redirect(w, r, "/", 301)
	}
	t := template.New("Authentification")
	t = template.Must(t.ParseFiles("./Client/test.html"))
	err := t.ExecuteTemplate(w, "test.html", ErrMessage)
	CheckError(err, "template")
}
