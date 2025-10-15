package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/dto"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/repository"
	filecreator "github.com/Jamaceat/liquibase-versioning-app/utils/file_creator"
)

func typeServiceHelper(repository repository.SchemaRepository,
	schema string,
	fatherContext context.Context,
	scriptUnformattedChan chan<- []string,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {

	defer wg.Done()

	typeNames, err := repository.GetTypesName(fatherContext, schema)

	if err != nil {
		errChan <- errors.New("not found typeNames")
		return
	}

	if err != nil {
		errChan <- errors.New("not found typeNames")
		return
	}

	var resultEntity []dto.TypesEntityComplete

	for _, value := range typeNames {
		var typesColumns []dto.TypeEntityColumns
		var typeComplete dto.TypesEntityComplete
		typeComplete.TypeName.TypName = value
		typeComplete.SchemaName = schema
		typesColumns, err = repository.GetTypesColumns(fatherContext, schema, value)
		if err != nil {
			errChan <- fmt.Errorf("error getting column data for type: %s", value)
			return
		}

		typeComplete.TypeEntityColumns = append(typeComplete.TypeEntityColumns, typesColumns...)

		resultEntity = append(resultEntity, typeComplete)
	}
	var scriptUnformatted []string
	for _, value := range resultEntity {

		scriptUnformatted = append(scriptUnformatted, filecreator.FormatType(value))
	}

	scriptUnformattedChan <- scriptUnformatted

}

func enumServiceHelper(repository repository.SchemaRepository,
	schema string,
	fatherContext context.Context,
	scriptUnformattedChan chan<- []string,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	enumsUnformatted, err := repository.GetEnums(fatherContext, schema)

	if err != nil {
		errChan <- fmt.Errorf("error al intentar recuperar los enums %w", err)
	}

}
