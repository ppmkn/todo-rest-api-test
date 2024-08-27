package main

import (
    "os"
	"log"
	"net/http"

	"todo-rest-api-test/internal/todo"

    "github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// инициализируем переменные подключения к БД
func init() {
	if err := godotenv.Load("../../configs/development.env"); err != nil {
		log.Fatal("Error loading .env file!")
	}
	todo.Port = os.Getenv("PORT")
	todo.DbDriver = os.Getenv("DB_DRIVER")
	todo.DbUser = os.Getenv("DB_USER")
	todo.DbPass = os.Getenv("DB_PASSWORD")
	todo.DbHost = os.Getenv("DB_HOST")
	todo.DbName = os.Getenv("DB_NAME")
	todo.DbMode = os.Getenv("DB_SSLMODE")
}

func main() {
	db, err := todo.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", todo.CreateTaskHandler(db)).Methods("POST")
	r.HandleFunc("/tasks", todo.GetTasksHandler(db)).Methods("GET")
	r.HandleFunc("/tasks/{id}", todo.GetTaskHandler(db)).Methods("GET")
	r.HandleFunc("/tasks/{id}", todo.UpdateTaskHandler(db)).Methods("PUT")
	r.HandleFunc("/tasks/{id}", todo.DeleteTaskHandler(db)).Methods("DELETE")

	log.Println("The server is running | Port:", todo.Port)
	log.Fatal(http.ListenAndServe(todo.Port, r))
}
