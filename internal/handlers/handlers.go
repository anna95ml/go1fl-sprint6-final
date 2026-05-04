package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
	//"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// В этом пакете вы реализуете два хендлера.
// Для корневого эндпоинта / нужно реализовать хендлер, который возвращает HTML из файла index.html.
// Второй хендлер для эндпоинта /upload должен выполнять следующие действия:
// Парсить html-форму из файла index.html.
// Получить файл из формы (не забудьте его закрыть).
// Прочитать данные из файла.
// Передать эти данные в функцию автоопределения из пакета service, которую вы создали, чтобы получить переконвертируемую строку.
// Создать локальный файл. Эта операция обычно небезопасна и так делать не рекомендуется, но в рамках нашего задания хотелось бы более наглядного результата, поэтому мы решились на этот шаг, ради видимого результата. А вообще, обычно используют временные файлы.
// Записать в локальный файл результат конвертации строки. Для генерации имени файла вы можете использовать время с помощью time.Now().UTC().String(). Чтобы получить расширения файла, используйте filepath.Ext().
// Вернуть результат конвертации строки.
// Там, где это необходимо, обработайте возможные ошибки. Статус при возникновении ошибок http.StatusInternalServerError.

func MainHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "../index.html") //сам запишет статус 200, если файл найден, или 404
}

func UploadHandle(res http.ResponseWriter, req *http.Request) {

	//проверяем метод
	if req.Method != http.MethodPost {
		http.Error(res, "метод отличается от POST", http.StatusMethodNotAllowed)
		return
	}

	// парсим html-форму
	err := req.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(res, "ошибка при парсинге формы"+err.Error(), http.StatusInternalServerError)
		return
	}

	// получаем файл из формы
	file, handler, err := req.FormFile("myFile")
	if err != nil {
		http.Error(res, "ошибка при получении файла из формы"+err.Error(), http.StatusInternalServerError)
		return
	}

	// закрываем файл
	defer file.Close()

	//	читаем данные из файла.
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(res, "ошибка при чтении файла", http.StatusInternalServerError)
		return
	}
	// конвертируем в морзе или обратно
	result, err := service.Translate(string(data))
	if err != nil {
		http.Error(res, "ошибке при переводе"+err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println("Текст из формы:", string(data), " Перевод:", result)

	//создаем локальный файл
	//получаем расширение файла
	ext := filepath.Ext(handler.Filename)
	// генерируем строку времени
	//timeStr := time.Now().UTC().String() //Отлетает в ошибку из-за двоеточий
	timeStr := time.Now().UTC().Format("20060102_150405.123")
	//формируем имя файла
	filename := handler.Filename + " " + timeStr + ext
	//fmt.Println(filename)

	//открываем директорию для загрузок
	root, err := os.OpenRoot("../uploads")
	if err != nil {
		http.Error(res, "внутренняя ошибка"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer root.Close()

	// создаём файл с таким же именем
	dst, err := root.Create(filename)
	if err != nil {
		http.Error(res, "ошибка при создании файла"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	//пишем в локальный файл результат конвертации строки.
	_, err = dst.WriteString(result)
	if err != nil {
		http.Error(res, "ошибка при записи файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//возвращаем результат конвертации строки.
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(result))
}
