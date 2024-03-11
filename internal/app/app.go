package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/LaQuannT/astronaut-data-api/internal/config"
	"github.com/LaQuannT/astronaut-data-api/internal/database"
	"github.com/LaQuannT/astronaut-data-api/internal/store"
	"github.com/LaQuannT/astronaut-data-api/internal/transport"
)

func Run() {
	env := config.Init()

	logger := config.InitLogger(os.Stdout, env.Stage)

	db := database.NewPostgresDB(env.BuildDBConnStr(), logger)
	dbPool, err := db.Init()
	if err != nil {
		logger.Log(context.Background(), config.LevelTrace, "failed database initialization", slog.Any("error", err))
		os.Exit(1)
	}

	us := store.NewUserStore(dbPool)
	as := store.NewAstronautStore(dbPool)

	addr := fmt.Sprintf(":%s", env.Port)

	s := transport.NewServer(addr, us, as, logger)
	s.Serve()
}
