package app

import (
	"fmt"
	"log"

	"github.com/LaQuannT/astronaut-data-api/internal/config"
	"github.com/LaQuannT/astronaut-data-api/internal/database"
	"github.com/LaQuannT/astronaut-data-api/internal/store"
	"github.com/LaQuannT/astronaut-data-api/internal/transport"
)

func Run() {
	env := config.Init()

	db := database.NewPostgresDB(env.BuildDBConnStr())
	dbPool, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	us := store.NewUserStore(dbPool)
	as := store.NewAstronautStore(dbPool)

	addr := fmt.Sprintf(":%s", env.Port)

	s := transport.NewServer(addr, us, as)
	s.Serve()
}
