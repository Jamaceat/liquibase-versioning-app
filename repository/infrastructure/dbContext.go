package context_db

import (
	"database/sql"
	"fmt"
	"log"

	constant "github.com/Jamaceat/liquibase-versioning-app/constants"
	"github.com/Jamaceat/liquibase-versioning-app/model/dto"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type DbContext struct {
	Database *sql.DB
}

var Db DbContext

// This method init the connection to the BD.
func init() {
	fmt.Println("Iniciando la BD.")

	envs, err := godotenv.Read(constant.ENV_FILE)
	if err != nil {
		panic("Cannot read .env file")
	}

	connectionConfiguration := dto.BDConnectionConfiguration{Host: envs[constant.DBHost],
		Port:     envs[constant.DBHost],
		User:     envs[constant.DBUser],
		Password: envs[constant.DBPassword],
		DBName:   envs[constant.DBName],
	}

	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", connectionConfiguration.User, connectionConfiguration.Password, connectionConfiguration.Host, connectionConfiguration.Port, connectionConfiguration.DBName))
	if err != nil {
		log.Fatal("No se pudo establecer conexión con la BD.")

	}

	fmt.Println("BD iniciada.")
	Db = DbContext{db}
}
