package transport

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	usecase "github.com/LaQuannT/astronaut-data-api/internal/service"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/handler"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/middleware"
	"github.com/gorilla/mux"
)

type server struct {
	log            *slog.Logger
	userStore      model.UserStore
	astronautStore model.AstronautStore
	addr           string
}

func NewServer(addr string, us model.UserStore, as model.AstronautStore, l *slog.Logger) *server {
	return &server{
		addr:           addr,
		userStore:      us,
		astronautStore: as,
		log:            l,
	}
}

func (s *server) Serve() {
	r := mux.NewRouter()

	sr := r.PathPrefix("/api/v1").Subrouter()
	sr.Use(middleware.HTTPLogger(s.log))

	userService := usecase.NewUserUsecase(s.userStore)
	astronautService := usecase.NewAstronautUsecase(s.astronautStore, s.userStore)

	handler.RegisterUserHandlers(userService, sr, s.log)
	handler.RegisterAstronautHandlers(astronautService, sr, s.log)

	s.log.Info(fmt.Sprintf("Server listening on '%s'", s.addr))
	log.Fatal(http.ListenAndServe(s.addr, r))
}
