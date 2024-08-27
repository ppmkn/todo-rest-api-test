package todo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // драйвер для postgres
)

// переменные для конфигурации БД
var (
	Port         string
	DbDriver     string
	DbUser       string
	DbPass       string
	DbHost       string
	DbName       string
	DbMode       string
	DbConnection string
)

// ConnectDatabase - реализует подключение к базе данных
func ConnectDatabase() (*sql.DB, error) {
	// формируем строку подключения
	DbConnection = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s", DbUser, DbPass, DbHost, DbName, DbMode)

	db, err := sql.Open(DbDriver, DbConnection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
