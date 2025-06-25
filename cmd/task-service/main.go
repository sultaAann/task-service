package main

import (
	"net/http"
	"task-sevice/internal/handler"
	"task-sevice/internal/repository"
	"task-sevice/internal/service"
)

func main() {
	db := repository.NewInstanceDB()

	repos := repository.NewRepository(db)

	service := service.NewService(repos)

	handler := handler.NewHandler(service)

	mux := http.NewServeMux()

	mux.Handle("/task/", handler)
	mux.Handle("/task", handler)

	http.ListenAndServe(":8080", mux)
}
