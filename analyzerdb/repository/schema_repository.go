package repository

import (
	"context"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/dto"
	context_db "github.com/Jamaceat/liquibase-versioning-app/repository/infrastructure"
)

type (
	SchemaRepository interface {
		GetExtensions(context context.Context) (extEntity []dto.ExtensionEntity, err error)
	}

	schemaRepository struct {
		db context_db.DbContext
	}
)

// GetExtensions implements SchemaRepository.
func (s *schemaRepository) GetExtensions(context context.Context) (extEntity []dto.ExtensionEntity, err error) {

	rows, err := s.db.Database.QueryContext(context, getExtensions)

	if err != nil {
		return extEntity, err
	}
	defer rows.Close()

	var extensionEntityAux dto.ExtensionEntity

	for rows.Next() {

		err = rows.Scan(
			&extensionEntityAux.Extname,
			&extensionEntityAux.Extversion,
		)

		if err != nil {
			break

		}
		extEntity = append(extEntity, extensionEntityAux)
	}

	return

}

func NewSchemaRepository() SchemaRepository {

	return &schemaRepository{context_db.Db}

}
