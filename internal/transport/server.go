package transport

import (
	"log"
	"net/http"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	usecase "github.com/LaQuannT/astronaut-data-api/internal/service"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/handler"
	"github.com/gorilla/mux"
)

type server struct {
	userStore      model.UserStore
	astronautStore model.AstronautStore
	addr           string
}

func NewServer(addr string, us model.UserStore, as model.AstronautStore) *server {
	return &server{
		addr:           addr,
		userStore:      us,
		astronautStore: as,
	}
}

func (s *server) Serve() {
	r := mux.NewRouter()

	sr := r.PathPrefix("/api/v1").Subrouter()

	userService := usecase.NewUserUsecase(s.userStore)
	astronautService := usecase.NewAstronautUsecase(s.astronautStore, s.userStore)

	handler.RegisterUserHandlers(userService, sr)
	handler.RegisterAstronautHandlers(astronautService, sr)

	log.Printf("Server listening on %q", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, r))
}
