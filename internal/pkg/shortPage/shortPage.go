package shortPage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonStruct struct {
	// Структура для декодирования и кодирования JSON
	Url string `json:"url"`
}

type UrlDBRecord struct {
	// Структура записей в базе данных
	Id        int64  `json:"id"`
	Long_url  string `json:"long_url"`
	Short_url string `json:"short_url"`
}

func ShortPage(w http.ResponseWriter, r *http.Request) {
	/* По условиям задачи при переходе на эту страницу на вход поступает длинная ссылка,
	возвращается сокращённая ссылка в формате JSON */
	var urlStruct JsonStruct
	errDecode := json.NewDecoder(r.Body).Decode(&urlStruct) // Декодирование поступишего на вход JSON
	if errDecode != nil {
		fmt.Println("ошибка в декодировании Json")
		panic(errDecode)
	}
	longUrl := urlStruct.Url
	shortUrl := getShortLink(longUrl)    // из этой функции получаем сокращенную ссылку
	newStructWithShortUrl := JsonStruct{ // создаем и заполняем структуру для последующего кодирования в JSON
		Url: shortUrl,
	}
	jsonInApiAnswer, errMar := json.Marshal(newStructWithShortUrl) // Кодирование в JSON
	if errMar != nil {
		fmt.Println("error Marshal Json")
		panic(errMar)
	}
	w.Header().Set("Content-Type", "application/json") // Правка заголовков и формирование тела ответа
	w.Write(jsonInApiAnswer)
}

func getShortLink(long string) string {
	/* Функция для сокращения ссылки. Вначале проверяем наличие в БД подобной записи если
	уже есть подобная запись то возвращаем ее, если нет то выполняем сокращение по алфавиту [A-Za-z0-9]
	основанием в нашем случаем является localhost:8000/{Запись основанная на полученном ID записи} */
	const (
		letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" // Алфавит нашей системы исчисления
		length  = int64(len(letters))
	)
	db, _ := sql.Open("sqlite3", "url_service.db")
	defer db.Close()
	row := db.QueryRow("SELECT * FROM urls WHERE long_url=?", long) //Проверяем существует ли уже у нас в базе данная ссылка
	answer := UrlDBRecord{}
	errRowScan := row.Scan(&answer.Id, &answer.Long_url, &answer.Short_url)
	if errRowScan == nil { // если такая запись уже есть то возвращаем эту запись
		return answer.Short_url
	}
	_, errInsert := db.Exec("INSERT INTO urls (long_url) VALUES (?)", long) // Если такой записи не существует то создаем запись в таблице с данной ссылкой
	if errInsert != nil {
		fmt.Println("Добавление записи не успешно")
		panic(errInsert)
	}
	linkProperties := UrlDBRecord{}
	rowId := db.QueryRow("SELECT id FROM urls WHERE long_url=?", long) // узнаем id который будем переводить в другую систему исчисления
	err := rowId.Scan(&linkProperties.Id)
	if err != nil {
		panic(err)
	}
	linkId := linkProperties.Id
	shortLink := "localhost:8000/"
	number := int64(linkId)                      // создаем переменную number чтобы переменная linkId не меняла своего значения
	for ; number > 0; number = number / length { // создаем короткую ссылку в новой системе исчисления
		shortLink = shortLink + string(letters[number%length])
	}
	_, err90 := db.Exec("UPDATE urls SET short_url=? WHERE id=?", shortLink, linkId)
	if err90 != nil {
		fmt.Println("Ошибка в UPDATE БД")
		panic(err90)
	}
	return shortLink // возвращаем итоговую сокращенную ссылку
}
