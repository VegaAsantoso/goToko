package app

import "github.com/VegaASantoso/goToko/app/controllers"

// Router
func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}