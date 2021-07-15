package longPage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonStruct struct { // Структура для декодирования и кодирования JSON
	Url string `json:"url"`
}

type UrlDBRecord struct { // Структура записей в базе данных
	Id        int64  `json:"id"`
	Long_url  string `json:"long_url"`
	Short_url string `json:"short_url"`
}

func longPage(w http.ResponseWriter, r *http.Request) {
	/* По условиям задачи при переходе на эту страницу на вход поступает длинная ссылка,
	возвращается сокращённая ссылка в формате JSON */
	var urlStruct JsonStruct
	errDecode := json.NewDecoder(r.Body).Decode(&urlStruct) // Декодирование поступишего на вход JSON
	if errDecode != nil {
		fmt.Println("ошибка в декодировании Json")
		panic(errDecode)
	}
	shortUrl := urlStruct.Url
	db, _ := sql.Open("sqlite3", "../../url_service.db")
	defer db.Close()
	// Проверка есть ли у нас в базе данных ссылки с такой сокращенной ссылкой
	row := db.QueryRow("SELECT * FROM urls WHERE short_url=?", shortUrl)
	answer := JsonStruct{}
	urlDB := UrlDBRecord{}
	err := row.Scan(&urlDB.Id, &urlDB.Long_url, &urlDB.Short_url)
	if err != nil {
		// Если нет такой ссылки то в формате JSON отправляем сообщение об отсутствии записи
		answer = JsonStruct{
			Url: "Такого URL в базе данных нет",
		}
	} else {
		// Если есть такая запись в БД то формируем ответ
		answer = JsonStruct{
			Url: urlDB.Long_url,
		}
	}
	jsonInApiAnswer, errMar := json.Marshal(answer) // Кодирование в JSON
	if errMar != nil {
		fmt.Println("error Marshal Json")
		panic(errMar)
	}
	w.Header().Set("Content-Type", "application/json") // Правка заголовков и формирование тела ответа
	w.Write(jsonInApiAnswer)

}
