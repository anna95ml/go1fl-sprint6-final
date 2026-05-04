package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Создайте структуру сервера с полями для логгера (log.Logger) и http-сервера (http.Server).
type Server struct {
	Lgr     *log.Logger
	HttpSrv *http.Server
}

// Создайте функцию, в которой нужно создать http-роутер. Функция принимает log.Logger и возвращает экземпляр структуры вашего сервера.
func CreateServer(logger *log.Logger) *Server {
	//fmt.Println("Создаем сервер из логгера")

	//Зарегистрируйте ваши хендлеры в http-роутере.
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.MainHandle)
	mux.HandleFunc(`/upload`, handlers.UploadHandle)
	//Создайте экземпляр структуры http.Server.
	s := &http.Server{
		Addr:         ":8080",          //Addr — используйте порт 8080.
		Handler:      mux,              //Handler — передайте ваш http-роутер.
		ErrorLog:     logger,           //ErrorLog — передайте ваш логгер.
		ReadTimeout:  5 * time.Second,  //ReadTimeout — таймаут для чтения. 5 секунд.
		WriteTimeout: 10 * time.Second, //WriteTimeout — таймаут для записи. 10 секунд.
		IdleTimeout:  15 * time.Second, //IdleTimeout — таймаут ожидания следующего запроса. 15 секунд.
	}
	return &Server{
		Lgr:     logger,
		HttpSrv: s,
	}
}
