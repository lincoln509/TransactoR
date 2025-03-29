package router

import (
	"net/http"

	"TransactoR/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Router struct {
	*mux.Router
	db *gorm.DB
}

func New(db *gorm.DB) *Router {
	r := mux.NewRouter()
	r.Use(middleware.Transaction(db))
	return &Router{Router: r, db: db}
}

func (r *Router) AddRoute(path, method string, handler http.HandlerFunc) {
	r.HandleFunc(path, handler).Methods(method)
}
