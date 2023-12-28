package routes

import(
	"github.com/gorilla/mux"
	"github.com/restingdemon/go-mysql-mux/pkg/controllers"
)

var RegisterAuthRoutes = func (router *mux.Router)  {
	router.HandleFunc("/signup",controllers.Signup).Methods("POST")
	router.HandleFunc("/login",controllers.Login).Methods("POST")
}