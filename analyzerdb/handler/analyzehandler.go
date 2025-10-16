package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/service"
	sliceutils "github.com/Jamaceat/liquibase-versioning-app/utils/slice_utils"
)

type (
	AnalyzerDB interface {
		GetDatabaseMigration() Controller
		GetTypeMigration() Controller
		GetTablesMigration() Controller
		GetParametersData() Controller
	}

	analyzerPostgres struct {
		servcExtension service.ExtensionService
	}
)

// GetParametersData implements AnalyzerDB.
func (a *analyzerPostgres) GetParametersData() Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		// buffer,err :=
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		if err := r.ParseMultipartForm((10 << 20) * 2); err != nil {
			http.Error(w, "El archivo es demasiado grande.", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("parameterizedTables")
		if err != nil {
			fmt.Println("Error al obtener el archivo", err)
			http.Error(w, "No se pudo obtener el archivo del formulario.", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "No se pudo leer el archivo.", http.StatusBadRequest)
			return
		}

		fileContent := string(fileBytes)
		lines := strings.Split(fileContent, "\n")

		uniques := sliceutils.UniquesStringFamily(lines)
		fmt.Println("---------------------------------------------")
		for _, v := range uniques {
			fmt.Println(v)
		}
		fmt.Println("---------------------------------------------")
		fmt.Printf("Archivo leido con exito %s", handler.Filename)

	}
}

// GetTablesMigration implements AnalyzerDB.
func (a *analyzerPostgres) GetTablesMigration() Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		buffer, err := a.servcExtension.GetFilesFromTables("public")

		if err != nil {
			log.Printf("Error generando el archivo de migracion TypeService en el buffer: %v", err)
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			return
		}
		downloadFilename := fmt.Sprintf("migration-buffered-tables-%s.sql", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
		w.Header().Set("Content-Length", strconv.Itoa(buffer.Len())) // .Len() da el tamaño en bytes.

		_, err = w.Write(buffer.Bytes()) // .Bytes() devuelve el contenido como un slice de bytes.
		if err != nil {
			log.Printf("Error escribiendo el buffer en la respuesta: %v", err)
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)

	}
}

// GetTypeMigration implements AnalyzerDB.
func (a *analyzerPostgres) GetTypeMigration() Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer, err := a.servcExtension.GetFilesFromTypes("public")

		if err != nil {
			log.Printf("Error generando el archivo de migracion TypeService en el buffer: %v", err)
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			return
		}
		downloadFilename := fmt.Sprintf("migration-buffered-types-%s.sql", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
		w.Header().Set("Content-Length", strconv.Itoa(buffer.Len())) // .Len() da el tamaño en bytes.

		_, err = w.Write(buffer.Bytes()) // .Bytes() devuelve el contenido como un slice de bytes.
		if err != nil {
			log.Printf("Error escribiendo el buffer en la respuesta: %v", err)
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)

	}
}

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
