package routes

import (
	"github.com/gorilla/mux"
	"github.com/rampl1ng/go-todoList/controller"
)

func Init() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.Show)
	r.HandleFunc("/add", controller.Add).Methods("POST")
	return r
}
