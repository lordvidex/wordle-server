package auth

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHTTPHandler(router *mux.Router) {
	router.HandleFunc("/login", LoginHTTPHandler).Methods("POST")
	// router.HandleFunc("/register", RegisterHTTPHandler).Methods("POST")
	// router.HandleFunc("/check", CheckHandler).Methods("POST")
}

func LoginHTTPHandler(http.ResponseWriter, *http.Request) {

}
