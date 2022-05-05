package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rampl1ng/go-todoList/config"
)

var (
	id        int
	item      string
	completed bool
	view      = template.Must(template.ParseFiles("./views/index.html"))
	db        = config.Database()
)

type View struct {
	Todos []Todo
}

type Todo struct {
	Id        int
	Item      string
	Completed bool
}

func Show(w http.ResponseWriter, r *http.Request) {
	todos := make([]Todo, 8)

	rows, err := db.Query(`SELECT * FROM todos`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &item, &completed)
		if err != nil {
			fmt.Println(err)
		}
		todo := Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
		}
		todos = append(todos, todo)
	}
	data := View{
		Todos: todos,
	}
	_ = view.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")
	_, err := db.Exec(`INSERT INTO todos (item) VALUE (?)`, item)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// todo: undo complete
func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
