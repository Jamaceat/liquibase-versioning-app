package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/dto"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/repository"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/templates"
	filecreator "github.com/Jamaceat/liquibase-versioning-app/utils/file_creator"
)

type (
	ExtensionService interface {
		GetFilesFromExtensions() (file bytes.Buffer, err error)
		GetFilesFromTypes(schema string) (file bytes.Buffer, err error)
		GetFilesFromTables(schema string) (file bytes.Buffer, err error)
		GetDMLFromTables(schema string) (file bytes.Buffer, err error)
	}

	extensionService struct {
		schemaDb repository.SchemaRepository
	}
)

// GetDMLFromTables implements ExtensionService.
func (e *extensionService) GetDMLFromTables(schema string) (file bytes.Buffer, err error) {
	panic("unimplemented")
}

// GetFilesFromTables implements ExtensionService.
func (e *extensionService) GetFilesFromTables(schema string) (file bytes.Buffer, err error) {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Hour)
	defer cancel()

	numWorkers := 20
	tableNameChan := make(chan dto.TableName)
	errChan := make(chan error, 1) // Buffered channel para no bloquear al primer error
	tableDetailedChan := make(chan dto.TableDetailed)

	var wg sync.WaitGroup

	// 1. Iniciar los workers
	for i := range numWorkers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for value := range tableNameChan {
				tableDetail, err := e.schemaDb.GetTableDetail(context, schema, value.TableName)
				if err != nil {
					// Intentar enviar el error sin bloquear. Si el canal está lleno, se descarta.
					// Mejor aún sería usar un errgroup.
					select {
					case errChan <- err:
					default:
					}
					return // Salir del worker si hay un error
				}
				tableDetailedChan <- dto.TableDetailed{TableName: value.TableName, SchemaName: value.SchemaName, Columns: tableDetail}
			}
			log.Printf("Worker %d finalizado", id)
		}(i)
	}

	// 2. Enviar trabajos en una goroutine separada para no bloquear la principal
	go func() {
		tableNames, err := e.schemaDb.GetTablesName(context, schema)
		if err != nil {
			errChan <- err
		} else {
			for _, tableName := range tableNames {
				tableNameChan <- tableName
			}
		}
		// ¡Muy importante! Cerrar el canal de trabajos para que los workers terminen.
		close(tableNameChan)
	}()

	// 3. Goroutine para cerrar los canales de resultados cuando los workers terminen
	go func() {
		wg.Wait()
		close(tableDetailedChan)
		close(errChan)
	}()

	// 4. Recolectar resultados en la goroutine principal
	var resultBuilder strings.Builder
	var firstErr error

	// Usamos un bucle for/range doble para leer los canales hasta que se cierren.
	// Esto es más limpio que el bucle infinito con select y checks de `nil`.
	for tableDetailedChan != nil || errChan != nil {
		select {
		case errVal, ok := <-errChan:
			if !ok {
				errChan = nil // El canal se cerró
			} else if firstErr == nil {
				firstErr = errVal // Capturamos solo el primer error
			}
		case tableDetailed, ok := <-tableDetailedChan:
			if !ok {
				tableDetailedChan = nil // El canal se cerró
			} else {
				resultBuilder.WriteString(filecreator.FormatTables(tableDetailed))
			}
		}
	}

	if firstErr != nil {
		return bytes.Buffer{}, firstErr // Devolvemos el primer error que encontramos
	}

	filecreator.GenerateSQLFile(&file, resultBuilder.String())

	return file, nil // El 'named return' `err` será nil aquí.

}

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
