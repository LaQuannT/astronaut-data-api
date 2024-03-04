package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LaQuannT/astronaut-data-api/internal/database"
	usecase "github.com/LaQuannT/astronaut-data-api/internal/service"
	"github.com/LaQuannT/astronaut-data-api/internal/store"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/handler"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Connecting to database...")
	dbPool, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	userStore := store.NewUserStore(dbPool)
	astroStore := store.NewAstronautStore(dbPool)

	userService := usecase.NewUserUsecase(userStore)
	astroService := usecase.NewAstronautUsecase(astroStore, userStore)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}).Methods("GET")

	us := r.PathPrefix("/users").Subrouter()
	as := r.PathPrefix("/astronauts").Subrouter()

	handler.RegisterUserHandlers(userService, us)
	handler.RegisterAstronautHandlers(astroService, as)

	fmt.Println("Server is now listing on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
