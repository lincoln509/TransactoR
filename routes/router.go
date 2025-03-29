package routes

import (
	"net/http"

	// "TransactoR/logging"
	// "TransactoR/middleware"

	// "lib/middleware"

	"github.com/gorilla/mux"
	logging "github.com/lincoln509/TransactoR/middleware"
	// "github.com/lincoln509/TransactoR/middleware"
)

type RouteConfig struct {
	Path        string
	Method      string
	Handler     http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
}

type Router struct {
	*mux.Router
	TransactionHandler *middleware.TransactionHandler
}

func NewRouter(logger logging.Logger) *Router {
	r := mux.NewRouter()
	th := middleware.NewTransactionHandler(logger)

	return &Router{
		Router:             r,
		TransactionHandler: th,
	}
}

func (r *Router) AddTransactionalRoute(config RouteConfig) {
	handlerChain := r.TransactionHandler.Wrap(config.Handler)

	for _, mw := range config.Middlewares {
		handlerChain = mw(handlerChain)
	}

	r.Handle(config.Path, handlerChain).Methods(config.Method)
}
