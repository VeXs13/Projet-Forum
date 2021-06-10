package pages

import (
	"net/http"

	models "../../models"
)

func Disconnect(w http.ResponseWriter, r *http.Request, DisconnectButton string) {
	if DisconnectButton == "" {
		return
	}
	CurrUser = models.User{}
	http.Redirect(w, r, "/login", 301)
}
