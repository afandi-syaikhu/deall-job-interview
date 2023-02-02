package main

import (
	"database/sql"
	"fmt"

	"github.com/afandi-syaikhu/deall-job-interview/config"
	"github.com/afandi-syaikhu/deall-job-interview/delivery/rest"
	"github.com/afandi-syaikhu/deall-job-interview/pkg"
	"github.com/afandi-syaikhu/deall-job-interview/repository"
	"github.com/afandi-syaikhu/deall-job-interview/usecase"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {

	// read config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// init db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// init repo
	userRepo := repository.NewUserRepository(db)

	// init usecase
	authUC := usecase.NewAuthUseCase(userRepo, cfg)

	// init echo framework
	e := echo.New()
	e.Validator = pkg.GetValidator()

	// init handler
	rest.NewAuthHandler(e, authUC)

	e.Logger.Fatal(e.Start(":8080"))
}
