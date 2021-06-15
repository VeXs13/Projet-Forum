package pages

import (
	"fmt"
	"html/template"
	"net/http"

	controllers ".."
	data "../../../Database"
)

func UpdateDescription(Description string) {
	statement, err := data.Database.Prepare("UPDATE Utilisateur SET Description = '" + Byte13To10(Description) + "' WHERE ID = " + fmt.Sprint(CurrUser.ID))
	CheckError(err, "Profil statement line 14")
	statement.Exec()
	CurrUser.Description = Description
}

func UpdatePassword(Last, New string) (bool, string) {
	rows, err := data.Database.Query("SELECT password FROM Utilisateur WHERE ID = " + fmt.Sprint(CurrUser.ID))
	CheckError(err, "line 18 rows")
	defer rows.Close()

	var Password string
	for rows.Next() {
		rows.Scan(&Password)
	}

	if controllers.CheckPasswordHash(Last, Password) {
		fmt.Println(New)
		New, err = controllers.HashPassword(New)
		fmt.Println(New)
		CheckError(err, "Hashing")
		statement, err := data.Database.Prepare("UPDATE Utilisateur SET PASSWORD = '" + New + "' WHERE ID = " + fmt.Sprint(CurrUser.ID))

		CheckError(err, "Profil statement line 30")
		_, err = statement.Exec()
		CheckError(err, "line 38 Profil Exec ")
		return true, ""
	} else {
		return false, "Mot de passe incorrect"
	}

}

func ChangeName(Name string) (bool, string) {
	rows, err := data.Database.Query("SELECT * FROM Utilisateur WHERE Pseudo='" + Name + "'")
	CheckError(err, "Rows")
	defer rows.Close()
	for rows.Next() {
		return false, "Le nom existe déjà"
	}

	statement, err := data.Database.Prepare("UPDATE Utilisateur SET Pseudo = '" + Name + "' WHERE Pseudo = '" + CurrUser.Pseudo + "'")
	CheckError(err, "new value line 27")
	_, err = statement.Exec()
	CheckError(err, "IMPORTANT ")
	statement, err = data.Database.Prepare("UPDATE Commentaires SET User = '" + Name + "' WHERE User = '" + CurrUser.Pseudo + "'")
	CheckError(err, "new value line 28")
	statement.Exec()
	statement, err = data.Database.Prepare("UPDATE Dislikes SET Autor = '" + Name + "' WHERE Autor = '" + CurrUser.Pseudo + "'")
	CheckError(err, "new value line 29")
	statement.Exec()
	statement, err = data.Database.Prepare("UPDATE Likes SET Autor = '" + Name + "' WHERE Autor = '" + CurrUser.Pseudo + "'")
	CheckError(err, "new value line 30")
	statement.Exec()
	statement, err = data.Database.Prepare("UPDATE Posts SET AutorName = '" + Name + "' WHERE AutorName = '" + CurrUser.Pseudo + "'")
	CheckError(err, "new value line 31")
	statement.Exec()
	CurrUser.Pseudo = Name

	return true, ""
}

func EvenlyChangeProfil(NewName, LastPassword, NewPassword, NewDescription string) {
	if NewName != "" {
		ok, Errmsg := ChangeName(NewName)
		if !ok {
			CurrUser.ErrMessage = Errmsg
			return
		}
	}

	if LastPassword != "" && NewPassword != "" {
		ok, Errmsg := UpdatePassword(LastPassword, NewPassword)
		if !ok {
			CurrUser.ErrMessage = Errmsg
			return
		}
	}

	if NewDescription != "" {
		UpdateDescription(NewDescription)
	}

	CurrUser.ErrMessage = ""
}

func ProfilPage(w http.ResponseWriter, r *http.Request) {
	EvenlyChangeProfil(r.FormValue("NewName"), r.FormValue("LastPassword"), r.FormValue("NewPassword"), r.FormValue("NewDescription"))
	// fmt.Fprintf(w, "ERREUR")
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./Client/Profil.html"))
	err := t.ExecuteTemplate(w, "Profil.html", CurrUser)
	CheckError(err, "template")
}
