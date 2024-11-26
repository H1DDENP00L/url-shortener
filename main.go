package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

var db *sql.DB

// генерация уникальной короткой ссылки через слайс на 6 байтоу
func generateShortLink() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortLink := make([]byte, 6)
	for i := range shortLink {
		shortLink[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortLink)
}

// проверка URL на валидность через регулярку
func isValidURL(url string) bool {
	// регулярка для адекватной структуры ссылки
	regex := `^(http|https)://([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(/[-a-zA-Z0-9@:%_+.~#?&//=]*)?$`
	re := regexp.MustCompile(regex)
	return re.MatchString(url)
}

// коннектимся к наше БДшке
func initDB() {
	var err error
	connStr := "postgres://postgres:zxc@localhost:5432/urlShortener?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
}

// функция для сохранения ссылки в БД
func saveLink(originalURL, shortLink string) error {
	_, err := db.Exec("INSERT INTO links (original_url, short_url) VALUES ($1, $2)", originalURL, shortLink)
	return err
}

// функция для получения оригинальной ссылки по короткой
func getOriginalURL(shortLink string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM links WHERE short_url = $1", shortLink).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

// проверка, существует ли уже короткая ссылка для данной оригинальной
func getShortLinkForOriginalURL(originalURL string) (string, error) {
	var shortLink string
	err := db.QueryRow("SELECT short_url FROM links WHERE original_url = $1", originalURL).Scan(&shortLink)
	if err != nil {
		return "", nil // Если записи нет, возвращаем пустую строку
	}
	return shortLink, nil
}

// ShortenURL принимает длинный URL и возвращает его сокращенную версию.
// Если URL уже существует в базе данных, возвращает его существующий короткий вариант.
// В случае ошибки вернет... ошибку)
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Делайте только POST запросы пожалуйста)... ", http.StatusMethodNotAllowed)
		return
	}

	// читаем тело запроса
	var reqBody struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// проверка валидности URL
	if !isValidURL(reqBody.URL) {
		http.Error(w, "Некорректный URL", http.StatusBadRequest)
		return
	}

	// проверяем, есть ли уже сокращенная версия этой ссылки
	shortLink, err := getShortLinkForOriginalURL(reqBody.URL)
	if err != nil {
		http.Error(w, "Не удалось проверить ссылку в базе данных", http.StatusInternalServerError)
		return
	}

	if shortLink == "" {
		// если ссылки еще нет в базе данных, генерируем новую короткую ссылку
		shortLink = generateShortLink()

		// сохраняем новую ссылку в БД
		err = saveLink(reqBody.URL, shortLink)
		if err != nil {
			http.Error(w, "Не удается сохранить ссылку", http.StatusInternalServerError)
			return
		}
	}

	// ответ с короткой ссылкой
	resp := map[string]string{"short_url": fmt.Sprintf("http://localhost:8080/%s", shortLink)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// обработчик редиректа по короткой ссылке
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Path[1:] // Получаем короткую ссылку из URL

	// Получаем оригинальную ссылку из БД
	originalURL, err := getOriginalURL(shortLink)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	// инициализируем коннект к БД
	initDB()
	// соединение закроется вместе с программой
	defer db.Close()

	// обработчики
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	// запускаем сервак
	port := "8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
