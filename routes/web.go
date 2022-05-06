package routes

import (
	"github.com/gorilla/mux"
	"github.com/rampl1ng/go-todoList/controller"
)

func Init() *mux.Router {
	r := mux.NewRouter()

	// todo-list/v1 use mysql
	r1 := r.PathPrefix("/v1").Subrouter()
	r1.HandleFunc("/", controller.Show)
	r1.HandleFunc("/add", controller.Add).Methods("POST")
	r1.HandleFunc("/complete/{id}", controller.Complete)
	r1.HandleFunc("/delete/{id}", controller.Delete)

	// todo-list/v2 use mongodb
	r2 := r.PathPrefix("/v2").Subrouter()
	r2.HandleFunc("/", controller.GetAll)
	r2.HandleFunc("/create", controller.Create).Methods("POST")
	r2.HandleFunc("/delete/{id}", controller.Remove)
	r2.HandleFunc("/update/{id}", controller.Update)
	return r
}
