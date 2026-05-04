package main

import (
	"log"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	//создаем логгер
	logger := log.Default()
	//fmt.Println("Создали логгер")

	srv := server.CreateServer(logger)
	//fmt.Println("Создали сервер из логгера")

	//fmt.Println("Запускаем сервер")
	logger.Printf("Запуск сервера %s", srv.HttpSrv.Addr)
	if err := srv.HttpSrv.ListenAndServe(); err != nil {
		logger.Fatalf("Ошибка сервера: %v", err)
		//srv.Lgr.Fatal(err)
	}

}
