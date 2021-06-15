package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	data "../../Database"
	models "../models"
)

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf(message+"execution : %s", err)
	}
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	CookieLastLogged := &http.Cookie{
		Name:  "LastLogged",
		Value: "",
	}
	CookieUser := &http.Cookie{
		Name:  "UserID",
		Value: "0",
	}

	http.SetCookie(w, CookieLastLogged)
	http.SetCookie(w, CookieUser)
}

func MaybeHaveCookie(User models.User, w http.ResponseWriter, r *http.Request) models.User {
	// function which check if the user have already a cookies, and then if they have an account, expired or not
	if User.Session != "" {
		return User
	}
	c, err := r.Cookie("UserID")
	if err != nil {
		return User
	}

	rows, err := data.Database.Query("SELECT * FROM Utilisateur WHERE ID = " + c.Value)
	CheckError(err, "Cookies statement")
	for rows.Next() {
		var nom string
		var prenom string
		var mail string
		var passwd string
		var adress string
		var uid string
		err = rows.Scan(&uid, &nom, &prenom, &mail, &passwd, &adress)
		CheckError(err, "Rows statement")
		User.ID = ConvertToInt(uid)
		User.Mail = mail
		User.Pseudo = prenom
		User.Session = "Connected"
	}
	return User
}

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	// Fonction qui va vérifier si la session dde l'utilisateur est expirée
	c, err := r.Cookie("LastLogged")
	fmt.Println(c)
	CheckError(err, "Cookie")
	timeint, _ := strconv.ParseInt(c.Value, 10, 0)
	LastActivtyTime, err := time.Parse("01-02-2006 15:04:05", time.Unix(timeint, 0).Format("01-02-2006 15:04:05"))
	fmt.Println(LastActivtyTime)
	CheckError(err, "Time")
	fmt.Println(time.Since(time.Now()), time.Now())
	TimeSinActivity := time.Since(LastActivtyTime)
	// fmt.Println(TimeSinActivity, TimeSinActivity.Minutes())
	// ! Fuseau Horaire
	if TimeSinActivity.Hours() > 3 { // If the Inactivity exeed x hours
		fmt.Println("Session expirated : ", TimeSinActivity.Hours())

		CookieUser := &http.Cookie{
			Name:  "UserID",
			Value: "0",
		}

		http.SetCookie(w, CookieUser)

		return false
	}

	return true

}

func ActualizeSession(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "LastLogged",
		Value: strconv.FormatInt(time.Now().Unix(), 10), // Recup the actual time
	}
	http.SetCookie(w, cookie) // Actualize the cookie with the new value
}
