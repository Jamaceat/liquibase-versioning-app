package handler

import (
	"net/http"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/service"
)

type (
	AnalyzerDB interface {
		GetDatabaseMigration(s service.ExtensionService) Controller
	}

	analyzerPostgres struct {
		servc service.ExtensionService
	}
)

// GetDatabaseMigration implements AnalyzerDB.
func (a *analyzerPostgres) GetDatabaseMigration(s service.ExtensionService) Controller {

	a.servc.GetFilesFromExtensions()

	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(200)

	}

}

func CreateEndpoint(servc service.ExtensionService) AnalyzerDB {

	return &analyzerPostgres{servc: servc}
}
