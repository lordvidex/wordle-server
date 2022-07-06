package auth

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler).Methods("POST")
	// router.HandleFunc("/register", RegisterHandler).Methods("POST")
	// router.HandleFunc("/check", CheckHandler).Methods("POST")
}

func LoginHandler(http.ResponseWriter, *http.Request) {

}

func LogoutHandler(http.ResponseWriter, *http.Request) {

}
