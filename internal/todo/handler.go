package todo

import (
    "encoding/json"
	"database/sql"
	"net/http"
	"strconv"
    "fmt"

	"github.com/gorilla/mux"
)

// CreateTaskHandler - обработчик создания задачи
func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Invalid JSON format!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		// дополнительная проверка
		if task.Title == "" || task.DueDate.IsZero() {
			http.Error(w, "Title and Due Date are required!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		createdTask, err := CreateTask(db, task.Title, task.Description, task.DueDate)
		if err != nil {
			http.Error(w, "Internal server error!", http.StatusInternalServerError) //500 Internal Server Error
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdTask)
	}
}

// GetTasksHandler - обработчик поиска всех задач
func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := GetTasks(db)
		if err != nil {
			http.Error(w, "Internal server error!", http.StatusInternalServerError) //500 Internal Server Error
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}
}

// GetTaskHandler - обработчик поиска конкретной задачи
func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid task ID!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		task, err := GetTask(db, id)
		if err != nil {
			http.Error(w, "Internal server error!", http.StatusInternalServerError) //500 Internal Server Error
			return
		}

		if task.ID == 0 {
			http.Error(w, "Task not found!", http.StatusNotFound) // 404 Not Found
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}
}

// UpdateTaskHandler - обработчик обновления конкретной задачи
func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid task ID!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		var task Task
		err = json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Invalid JSON format!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		// дополнительная проверка
		if task.Title == "" || task.DueDate.IsZero() {
			http.Error(w, "Title and Due Date are required!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		updatedTask, err := UpdateTask(db, id, task.Title, task.Description, task.DueDate)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Task not found!", http.StatusNotFound) //404 Not Found
				return
			}
			http.Error(w, "Internal server error!", http.StatusInternalServerError) // 500 Internal Server Error
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(updatedTask)
	}
}

// DeleteTaskHandler - обработчик удаления конкретной задачи
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid task ID!", http.StatusBadRequest) // 400 Bad Request
			return
		}

		err = DeleteTask(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Task not found!", http.StatusNotFound) // 404 Not Found
				return
			}
			http.Error(w, "Internal server error!", http.StatusInternalServerError) // 500 Internal Server Error
			return
		}

        fmt.Fprintf(w, "Task deleted successfully!")
		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}
}
