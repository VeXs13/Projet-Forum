package pages

import (
	"net/http"

	controllers ".."
	models "../../models"
)

func Disconnect(w http.ResponseWriter, r *http.Request, DisconnectButton string) {
	if DisconnectButton == "" {
		return
	}
	CurrUser = models.User{}
	controllers.DeleteCookie(w, r)
	http.Redirect(w, r, "/login", 301)
}
