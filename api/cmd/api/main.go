package main

import (
	"api/internal/handler"
	"api/internal/repository"
	"api/internal/service"
	"log"
	"net/http"
)

func main() {
	userRepo := repository.NewInMemoryUserRepo()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	r := handler.NewRouter(userHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}
