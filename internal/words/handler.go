package words

import (
	"github.com/lordvidex/wordle-wf/internal/common/werr"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func RegisterHandler(router *mux.Router, service Service) *Handler {
	h := &Handler{service}
	router.HandleFunc("", h.GetRandomWordHandler).Methods("GET")
	router.HandleFunc("/fail", h.FailWordHandler).Methods("GET")
	router.HandleFunc("/{id:[\\d]+}", h.GetWordWithIdHandler).Methods("GET")

	return h
}

func (h *Handler) GetRandomWordHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	word := h.service.GetRandomWord()
	word.WriteJSON(w)
}

func (h *Handler) GetWordWithIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists {
		panic(werr.BadRequest("id is required"))
	}
	intid, err := strconv.Atoi(id)
	if err != nil {
		panic(werr.BadRequest("id must be an integer"))
	}
	word, err := h.service.GetWordWithID(intid)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	word.WriteJSON(w)
}

func (h *Handler) FailWordHandler(w http.ResponseWriter, _ *http.Request) {
	panic(werr.InternalServerError("Yay yay"))
}
