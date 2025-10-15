package repository

import (
	"context"
	"fmt"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/dto"
	context_db "github.com/Jamaceat/liquibase-versioning-app/repository/infrastructure"
)

type (
	SchemaRepository interface {
		GetExtensions(context context.Context) (extEntity []dto.ExtensionEntity, err error)
		GetTypesName(context context.Context, schema string) (typesName []string, err error)
		GetTypesColumns(context context.Context, schema, typeName string) (typesEntityColumn []dto.TypeEntityColumns, err error)
		GetEnums(context context.Context, schema string) (enumEntity []dto.EnumEntity, err error)
	}

	schemaRepository struct {
		db context_db.DbContext
	}
)

// GetEnums implements SchemaRepository.
func (s *schemaRepository) GetEnums(context context.Context, schema string) (enumsEntity []dto.EnumEntity, err error) {
	rows, err := s.db.Database.QueryContext(context, fmt.Sprintf(getEnumDetailed, schema))

	if err != nil {
		return enumsEntity, err
	}
	defer rows.Close()

	var enumEntity dto.EnumEntity

	for rows.Next() {

		err = rows.Scan(
			&enumEntity.SchemaName,
			&enumEntity.EnumName,
			&enumEntity.EnumValue,
		)

		if err != nil {
			break

		}
		enumsEntity = append(enumsEntity, enumEntity)
	}

	return
}

// GetTypesDetailed implements SchemaRepository.
func (s *schemaRepository) GetTypesColumns(context context.Context, schema, typeName string) (typesEntityColumn []dto.TypeEntityColumns, err error) {

	rows, err := s.db.Database.QueryContext(context, getTypesDetailed)

	if err != nil {
		return typesEntityColumn, err
	}
	defer rows.Close()

	var typeEntityColumn dto.TypeEntityColumns

	for rows.Next() {

		err = rows.Scan(
			&typeEntityColumn.ColumnName,
			&typeEntityColumn.AliasType,
		)

		if err != nil {
			break

		}
		typesEntityColumn = append(typesEntityColumn, typeEntityColumn)
	}

	return

}

// GetTypesName implements SchemaRepository.
func (s *schemaRepository) GetTypesName(context context.Context, schema string) (typesName []string, err error) {
	rows, err := s.db.Database.QueryContext(context, fmt.Sprintf(getTypes, schema))

	if err != nil {
		return typesName, err
	}
	defer rows.Close()

	var TypeEntityNAux string

	for rows.Next() {

		err = rows.Scan(
			&TypeEntityNAux,
		)

		if err != nil {
			break

		}
		typesName = append(typesName, TypeEntityNAux)
	}

	return
}

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
