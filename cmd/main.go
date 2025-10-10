package main

import (
	"encoding/json"
	"log"
	"net/http"

	constant "github.com/Jamaceat/liquibase-versioning-app/constants"
	"github.com/Jamaceat/liquibase-versioning-app/handler/middleware"
	"github.com/Jamaceat/liquibase-versioning-app/server"
	"github.com/joho/godotenv"
)

func main() {

	// router := chi.NewRouter()

	// router.Use(middleware.ContentTypeMiddleware)

	serverConfiguration := getEnviroments()

	srv := server.New(serverConfiguration)

	router := srv.NewRouter()

	router.Use(middleware.ContentTypeMiddleware)

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {

		data := map[string]any{"hello": "world"}
		jsonData, err := json.Marshal(data)

		if err != nil {
			http.Error(w, "Error al generar el JSON", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	})

	if err := srv.Start(router); err != nil {
		log.Fatal(err)
	}

}

func getEnviroments() (serverConfiguration server.Options) {
	envs, err := godotenv.Read(constant.ENV_FILE)
	if err != nil {
		panic("Cannot read .env file")
	}

	serverConfiguration = server.Options{
		Identifier: envs[constant.AppIdentifier],
		Port:       envs[constant.AppPort],
		Host:       envs[constant.AppHost],
	}

	return
}
