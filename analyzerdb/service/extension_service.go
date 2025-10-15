package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/repository"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/templates"
	filecreator "github.com/Jamaceat/liquibase-versioning-app/utils/file_creator"
)

type (
	ExtensionService interface {
		GetFilesFromExtensions() (file bytes.Buffer, err error)
		GetFilesFromTypes(schema string) (file bytes.Buffer, err error)
	}

	extensionService struct {
		schemaDb repository.SchemaRepository
	}
)

// GetFilesFromTypes implements ExtensionService.
func (e *extensionService) GetFilesFromTypes(schema string) (file bytes.Buffer, err error) {
	contextTypes, cancel := context.WithTimeout(context.Background(), 30*time.Hour)
	defer cancel()

	var buffer bytes.Buffer

	var wg sync.WaitGroup

	var resultUnformattedTypes []string

	var firstError error

	wg.Add(2)
	resultChan := make(chan []string)
	errChan := make(chan error)

	go typeServiceHelper(e.schemaDb, schema, contextTypes, resultChan, errChan, &wg)
	go enumServiceHelper(e.schemaDb, schema, contextTypes, resultChan, errChan, &wg)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				resultUnformattedTypes = append(resultUnformattedTypes, result...)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else if firstError == nil {
				firstError = err
			}
		case <-contextTypes.Done():
			return buffer, fmt.Errorf("error timeout superado %w", contextTypes.Err())
		}

		if resultChan == nil && errChan == nil {
			break
		}

	}

	// resultUnformattedTypes = append(resultUnformattedTypes, resultUnformattedEnums...)
	if firstError != nil {

		return file, firstError
	}
	var result string
	for _, value := range resultUnformattedTypes {

		result += value

	}
	filecreator.GenerateSQLFile(&buffer, result)
	return buffer, nil
}

// GetFilesFromExtensions implements ExtensionService.
func (e *extensionService) GetFilesFromExtensions() (file bytes.Buffer, err error) {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var buffer bytes.Buffer

	extensions, err := e.schemaDb.GetExtensions(context)
	var script string
	for _, value := range extensions {
		script += fmt.Sprintf(templates.ExtensionTemplate, value.Extname, value.Extversion, value.Extname, value.Extversion)
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
