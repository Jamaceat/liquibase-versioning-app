package dbContext

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Jamaceat/liquibase-versioning-app/model/dto"
)

type DbContext struct {
}

// This method init the connection to the BD.
func (dbContext *DbContext) Init(connectionConfiguration dto.BDConnectionConfiguration) *sql.DB {
	fmt.Println("Iniciando la BD.")

	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", connectionConfiguration.User, connectionConfiguration.Password, connectionConfiguration.Host, connectionConfiguration.Port, connectionConfiguration.DBName))
	if err != nil {
		log.Fatal("No se pudo establecer conexión con la BD.")
		return nil
	}

	fmt.Println("BD iniciada.")
	return db
}
