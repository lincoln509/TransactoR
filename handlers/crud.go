package handlers

import (
	"TransactoR/router"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type CRUDHandler struct {
	Repo *Repository
}

func NewCRUDHandler(repo *Repository) *CRUDHandler {
	return &CRUDHandler{Repo: repo}
}

func (h *CRUDHandler) Create(w http.ResponseWriter, r *http.Request) {
	modelType := reflect.TypeOf(h.Repo.Model)
	newModel := reflect.New(modelType).Interface()

	if err := json.NewDecoder(r.Body).Decode(newModel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(newModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newModel)
}

func (h *CRUDHandler) Get(w http.ResponseWriter, r *http.Request) {
	// router.GetParam()
	identifier := router.GetParam(r, "id")
	// result, err := h.Repo.GetByIDOrUsername(identifier)
	result, err := h.Repo.Get(identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *CRUDHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (h *CRUDHandler) Update(w http.ResponseWriter, r *http.Request) {
	identifier := router.GetParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.Update(identifier, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CRUDHandler) Delete(w http.ResponseWriter, r *http.Request) {
	identifier := router.GetParam(r, "id")

	if err := h.Repo.Delete(identifier); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RCRUDRoutes(r *mux.Router, path string, handler *CRUDHandler) {
	// Collection routes
	r.HandleFunc(path, handler.Create).Methods("POST")
	r.HandleFunc(path, handler.GetAll).Methods("GET")

	// Item routes
	itemPath := path + "/{id}"
	r.HandleFunc(itemPath, handler.Get).Methods("GET")
	r.HandleFunc(itemPath, handler.Update).Methods("PUT", "PATCH")
	r.HandleFunc(itemPath, handler.Delete).Methods("DELETE")
}
