package todo

import (
	"database/sql"
	"time"
)

// CreateTask - создать задачу
func CreateTask(db *sql.DB, title, description string, dueDate time.Time) (Task, error) {
	var task Task
	query := `INSERT INTO tasks (title, description, due_date) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, title, description, dueDate).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return task, err
	}
	task.Title = title
	task.Description = description
	task.DueDate = dueDate
	return task, nil
}

// GetTasks - найти все задачи, обращение в БД
func GetTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, description, due_date, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTask - найти конкретную задачу, обращение в БД
func GetTask(db *sql.DB, id int) (Task, error) {
	var task Task
	query := "SELECT id, title, description, due_date, created_at, updated_at FROM tasks WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, nil // задача не найдена
		}
		return task, err
	}
	return task, nil
}

// UpdateTask - обновить конкретную задачу, обращение в БД
func UpdateTask(db *sql.DB, id int, title, description string, dueDate time.Time) (Task, error) {
	var task Task
	query := `
        UPDATE tasks SET title = $1, description = $2, due_date = $3 
        WHERE id = $4 
        RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, title, description, dueDate, id).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return task, err
	}
	task.Title = title
	task.Description = description
	task.DueDate = dueDate
	return task, nil
}

// DeleteTask - удалить конкретную задачу, обращение в БД
func DeleteTask(db *sql.DB, id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
