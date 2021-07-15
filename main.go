package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/irvasikov/url_service/internal/pkg/longPage"
	"github.com/irvasikov/url_service/internal/pkg/shortPage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./url_service.db")
	// создаем базу данных и формируем ее структуру
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS urls(
			id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			long_url TEXT,
			short_url TEXT
			)`)
	if err != nil {
		panic(err)
	}
	database.Close()
	http.HandleFunc("/short/", shortPage.ShortPage) // прослушиваем этот адрес для получения сокращенной ссылки
	http.HandleFunc("/long/", longPage.LongPage)    // прослушиваем этот адрес для получения длинной ссылки
	log.Fatal(http.ListenAndServe(":8000", nil))
}
