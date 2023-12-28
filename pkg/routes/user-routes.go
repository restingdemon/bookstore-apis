package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/go-mysql-mux/pkg/controllers"
	"github.com/restingdemon/go-mysql-mux/pkg/middleware"
)
var RegisterUserRoutes = func (router *mux.Router)  {

	router.Use(middleware.Authenticate)
	router.HandleFunc("/users",controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}",controllers.GetUser).Methods("GET")
}	