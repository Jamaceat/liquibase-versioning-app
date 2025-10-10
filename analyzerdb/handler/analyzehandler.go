package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/service"
)

type (
	AnalyzerDB interface {
		GetDatabaseMigration() Controller
	}

	analyzerPostgres struct {
		servcExtension service.ExtensionService
	}
)

// GetDatabaseMigration implements AnalyzerDB.
func (a *analyzerPostgres) GetDatabaseMigration() Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		buffer, err := a.servcExtension.GetFilesFromExtensions()

		if err != nil {
			log.Printf("Error generando el migracion extensionService en el buffer: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		downloadFilename := fmt.Sprintf("migration-buffered-%s.sql", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
		w.Header().Set("Content-Length", strconv.Itoa(buffer.Len())) // .Len() da el tamaño en bytes.

		_, err = w.Write(buffer.Bytes()) // .Bytes() devuelve el contenido como un slice de bytes.
		if err != nil {
			log.Printf("Error escribiendo el buffer en la respuesta: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)

	}

}

func CreateEndpoint(servcExtension service.ExtensionService) AnalyzerDB {

	return &analyzerPostgres{servcExtension: servcExtension}
}
