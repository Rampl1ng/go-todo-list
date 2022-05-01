package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

var (
	rnd *renderer.Render
	db  *mgo.Database
)

const (
	hostName       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type todoModel struct {
	ID        bson.ObjectId
	Title     string
	Completed bool
	CreatedAt time.Time
}

type todo struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
}

func main() {

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)

	r.Mount("/todo", todoHandlers())

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port ", port)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	if err != nil {
		log.Fatal(err)
	}
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func todoHandlers() http.Handler {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/", fetchTodo)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return r
}
