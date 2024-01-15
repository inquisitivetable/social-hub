package main

import (
	"SocialNetworkRestApi/api/internal/server/handlers"
	"SocialNetworkRestApi/api/internal/server/router"
	"SocialNetworkRestApi/api/pkg/db/seed"
	database "SocialNetworkRestApi/api/pkg/db/sqlite"
	"SocialNetworkRestApi/api/pkg/models"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	port int
}

func main() {
	config := &Config{port: 8000}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	//DATABASE
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.RunMigrateScripts(db)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("successfully migrated DB..")

	repos := models.InitRepositories(db)

	app := handlers.InitApp(repos, logger)

	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "seed":
			seed.Seed(db, logger, repos)
		default:
			break
		}

	}

	r := router.New(app)

	logger.Printf("Starting server on port %d\n", config.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.port), r); err != nil {
		logger.Fatal(err)
	}

}
