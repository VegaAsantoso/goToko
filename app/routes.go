package app

import (
	"github.com/VegaASantoso/goToko/app/controllers"
	"github.com/gorilla/mux"
)

// Router
func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}