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
		GetTablesName(context context.Context, schema string) (tableNames []dto.TableName, err error)
		GetTableDetail(context context.Context, schema, tableName string) (tableColumns []dto.Column, err error)
	}

	schemaRepository struct {
		db context_db.DbContext
	}
)

// GetTableDetail implements SchemaRepository.
func (s *schemaRepository) GetTableDetail(context context.Context, schema string, tableName string) (tableColumns []dto.Column, err error) {

	rows, err := s.db.Database.QueryContext(context, getTableDetail, schema, tableName)

	if err != nil {

		return
	}

	var tableColumn dto.Column

	for rows.Next() {

		err = rows.Scan(
			&tableColumn.ColumnName,
			&tableColumn.NoAliasType,
			&tableColumn.AliasType,
		)

		if err != nil {
			break
		}

		tableColumns = append(tableColumns, tableColumn)

	}

	return
}

// GetTablesName implements SchemaRepository.
func (s *schemaRepository) GetTablesName(context context.Context, schema string) (tableNames []dto.TableName, err error) {

	rows, err := s.db.Database.QueryContext(context, getTables, schema)

	if err != nil {

		return
	}

	var tableNameAux dto.TableName

	for rows.Next() {

		err = rows.Scan(
			&tableNameAux.SchemaName,
			&tableNameAux.TableName,
		)

		if err != nil {

			break
		}

		tableNames = append(tableNames, tableNameAux)

	}

	return
}

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

	rows, err := s.db.Database.QueryContext(context, fmt.Sprintf(getTypesDetailed, typeName, schema))

	if err != nil {
		return typesEntityColumn, err
	}
	defer rows.Close()

	var typeEntityColumn dto.TypeEntityColumns

	for rows.Next() {

		err = rows.Scan(
			&typeEntityColumn.ColumnName,
			&typeEntityColumn.NoAlias,
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

	var TypeEntityNAux dto.TypesEntityN

	for rows.Next() {

		err = rows.Scan(
			&TypeEntityNAux.TypName,
		)

		if err != nil {
			break

		}
		typesName = append(typesName, TypeEntityNAux.TypName)
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
