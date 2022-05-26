package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampl1ng/go-todoList/config"
	"github.com/rampl1ng/go-todoList/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// id        int
	// item      string
	// completed bool
	// view            = template.Must(template.ParseFiles("./views/index.html"))
	view = template.Must(template.ParseFiles("./views/index.html"))
	// db              = config.Database()
	mongoCollection = config.GetCollection()
)

// Mysql todo view
// type View struct {
// 	Todos []Todo
// }

// // Mysql todo
// type Todo struct {
// 	Id        int
// 	Item      string
// 	Completed bool
// }

// func Show(w http.ResponseWriter, r *http.Request) {
// 	todos := make([]Todo, 0)

// 	rows, err := db.Query(`SELECT * FROM todos`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	for rows.Next() {
// 		var todo Todo

// 		err = rows.Scan(&todo.Id, &todo.Item, &todo.Completed)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		todos = append(todos, todo)
// 	}
// 	data := View{
// 		Todos: todos,
// 	}
// 	_ = view.Execute(w, data)
// }

// // TODO: if repeat item, add fails
// func Add(w http.ResponseWriter, r *http.Request) {
// 	item := r.FormValue("item")
// 	_, err := db.Exec(`INSERT INTO todos (item) VALUE (?)`, item)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	http.Redirect(w, r, "/v1/", http.StatusMovedPermanently)
// }

// // Complete changes the todo status to Completed
// func Complete(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]
// 	_, err := db.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	http.Redirect(w, r, "/v1/", http.StatusMovedPermanently)
// }

// func Delete(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]
// 	_, err := db.Exec(`DELETE FROM todos WHERE id = ?`, id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	http.Redirect(w, r, "/v1/", http.StatusMovedPermanently)
// }

/**
  ------------------
 |                  |
 |  v2 Use MongoDB  |
 |                  |
  ------------------
**/

type MongoDBView struct {
	Todos []MongoDBTodo
}

type MongoDBTodo struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Item      string             `bson:"title,omitempty" json:"title,omitempty"`
	Completed bool               `bson:"complete" json:"complete"`
}

func GetAllTodos(c *gin.Context) {
	var todos []MongoDBTodo

	ctx, cancel := utils.TodoContext()
	defer cancel()

	filter := bson.M{}
	findOptions := options.Find()
	cursor, err := mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println(err)
	}
	for cursor.Next(ctx) {
		var todo MongoDBTodo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	data := MongoDBView{
		Todos: todos,
	}
	// fmt.Printf("%#v\n", data)

	_ = view.Execute(c.Writer, data)
}

func AddOneToDo(c *gin.Context) {
	var todo MongoDBTodo

	item := c.PostForm("item")
	todo = MongoDBTodo{
		Item:      item,
		Completed: false,
	}

	res, err := mongoCollection.InsertOne(context.TODO(), todo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.InsertedID)

	c.Redirect(http.StatusMovedPermanently, "/v1/")
}

func DeleteOneToDo(c *gin.Context) {
	id := c.Param("id")
	objId, err := convertObjectID(id)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := utils.TodoContext()
	defer cancel()

	_, err = mongoCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		fmt.Println(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/v1/")
}

// Update changes the todo status to Completed
func UpdateOneToDo(c *gin.Context) {
	id := c.Param("id")
	objId, err := convertObjectID(id)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := utils.TodoContext()
	defer cancel()

	// change complete status false to true
	update := bson.M{
		"$set": bson.M{"complete": true},
	}
	fmt.Println(update)

	_, err = mongoCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		fmt.Println(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/v1/")
}

// convert id in url to ObjectID
// e.g.
// 62752727979a6f62a19514bf -> ObjectID("62752727979a6f62a19514bf")
func convertObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id[10:34])
}
