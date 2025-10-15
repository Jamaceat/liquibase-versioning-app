package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/handler/middleware"
	"github.com/go-chi/chi"
)

type Options struct {
	Identifier string
	Port       string
	Host       string
}

// NewRouter implements Server.
func (o *Options) NewRouter() Router {
	router := chi.NewRouter()

	router.Use(
		middleware.ContentTypeMiddleware,
	)

	return &engine{
		chi:    router,
		prefix: o.Identifier,
	}
}

// Start implements Server.
func (o *Options) Start(router http.Handler) error {

	fmt.Println("Iniciando Servidor en http://127.0.0.1:8008")
	srv := &http.Server{
		Handler: http.TimeoutHandler(router, time.Hour*20, "TimeOut!"),
		Addr:    fmt.Sprintf("%s:%s", o.Host, o.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	return nil
}

type Server interface {
	NewRouter() Router

	Start(handler http.Handler) error
}

func New(opts Options) Server {

	// Agregar aqui configuraciones

	return &opts
}
