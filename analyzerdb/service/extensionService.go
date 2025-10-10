package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/repository"
	"github.com/Jamaceat/liquibase-versioning-app/utils/path"
)

type (
	ExtensionService interface {
		GetFilesFromExtensions() (filePath string, err error)
	}

	extensionService struct {
		schemaDb repository.SchemaRepository
	}
)

// GetFilesFromExtensions implements ExtensionService.
func (e *extensionService) GetFilesFromExtensions() (filePath string, err error) {
	var sb strings.Builder

	sb.WriteString("-- Prueba")

	directoryPath := fmt.Sprintf("%s/generatedFiles", path.GetFindProjectRoot())

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {

		os.Mkdir(directoryPath, os.ModePerm)
	}

	filePath = fmt.Sprintf("%s/generatedFiles/migration-%s-%d.sql", directoryPath, "extensions", time.Now().UnixNano())

	file, err := os.Create(filePath)

	if err != nil {
		return "", err

	}
	defer file.Close()

	_, err = file.WriteString(sb.String())

	if err != nil {
		return "", err
	}

	return

}

func NewExtensionService(repository repository.SchemaRepository) ExtensionService {

	return &extensionService{
		schemaDb: repository,
	}

}
