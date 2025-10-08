package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	constant "github.com/Jamaceat/liquibase-versioning-app/constants"
	"github.com/Jamaceat/liquibase-versioning-app/model/dto"
	"github.com/Jamaceat/liquibase-versioning-app/router/middleware"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {

	router := chi.NewRouter()

	router.Use(middleware.ContentTypeMiddleware)

	enviroments := getEnviroments()

	fmt.Printf("%+v\n", enviroments)

	fmt.Println("Iniciando Servidor en http://127.0.0.1:8008")

	srv := &http.Server{
		Handler: http.TimeoutHandler(router, time.Second*20, "TimeOut!"),
		Addr:    "127.0.0.1:8000",
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	// router := mux.NewRouter()

	// fmt.Println("Iniciando Sv")

	// srv := &http.Server{
	// 	Handler: http.TimeoutHandler(router, time.Second*20, "TimeOut!"),
	// 	Addr:    "127.0.0.1:8008",
	// }

	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Fatal(err)
	// }

}

func getEnviroments() (configuration dto.BDConnectionConfiguration) {
	envs, err := godotenv.Read(constant.ENV_FILE)
	if err != nil {
		panic("Cannot read .env file")
	}

	configuration = dto.BDConnectionConfiguration{Host: envs[constant.DBHost],
		Port:     envs[constant.DBHost],
		User:     envs[constant.DBUser],
		Password: envs[constant.DBPassword],
		DBName:   envs[constant.DBName],
	}

	return
}
