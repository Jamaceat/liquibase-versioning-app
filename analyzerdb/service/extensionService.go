package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/repository"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/templates"
	filecreator "github.com/Jamaceat/liquibase-versioning-app/utils/file_creator"
)

type (
	ExtensionService interface {
		GetFilesFromExtensions() (filePath bytes.Buffer, err error)
	}

	extensionService struct {
		schemaDb repository.SchemaRepository
	}
)

// GetFilesFromExtensions implements ExtensionService.
func (e *extensionService) GetFilesFromExtensions() (filePath bytes.Buffer, err error) {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var buffer bytes.Buffer

	extensions, err := e.schemaDb.GetExtensions(context)
	var script string
	for _, value := range extensions {
		script += fmt.Sprintf(templates.ExtensionTemplate, value.Extname, value.Extversion, value.Extname)
	}

	if err != nil {
		log.Println("Error al recuperar las extentions")
		return buffer, err

	}

	err = filecreator.GenerateSQLFile(&buffer, script)
	if err != nil {

		return buffer, err

	}

	return buffer, nil

}

func NewExtensionService(repository repository.SchemaRepository) ExtensionService {

	return &extensionService{
		schemaDb: repository,
	}

}
