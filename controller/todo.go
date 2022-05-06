package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rampl1ng/go-todoList/config"
	"github.com/rampl1ng/go-todoList/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

var (
	// id        int
	// item      string
	// completed bool
	view            = template.Must(template.ParseFiles("./views/index.html"))
	db              = config.Database()
	mongoDb         = config.MongoClient()
	mongoCollection = config.GetCollection()
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
	todos := make([]Todo, 0)

	rows, err := db.Query(`SELECT * FROM todos`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var todo Todo

		err = rows.Scan(&todo.Id, &todo.Item, &todo.Completed)
		if err != nil {
			fmt.Println(err)
		}
		todos = append(todos, todo)
	}
	data := View{
		Todos: todos,
	}
	_ = view.Execute(w, data)
}

// TODO: if repeat item, add fails
func Add(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")
	_, err := db.Exec(`INSERT INTO todos (item) VALUE (?)`, item)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Complete changes the todo status to Completed
func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

/**
  ------------------
 |                  |
 |  v2 Use MongoDB  |
 |                  |
  ------------------
**/
func GetAll(w http.ResponseWriter, r *http.Request) {
	var todos []Todo

	ctx, cancel := utils.TodoContext()
	defer cancel()

	filter := bson.M{}
	findOptions := options.Find()
	cursor, err := mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println(err)
	}
	for cursor.Next(ctx) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	data := View{
		Todos: todos,
	}
	_ = view.Execute(w, data)
}

// TODO: test conncetion?
func chooseDB() {
	if err := mongoDb.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
	}
}
