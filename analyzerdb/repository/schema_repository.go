package repository

import (
	"context"
	"fmt"
	"regexp"

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
		GetDataForParameterizedTable(ctx context.Context, schema, tableName string) (rowsResult []map[string]any, err error)
	}

	schemaRepository struct {
		db context_db.DbContext
	}
)

var safeIdentifierRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// GetDMLForParameterizedTable implements SchemaRepository.
func (s *schemaRepository) GetDataForParameterizedTable(ctx context.Context, schema string, tableName string) (rowsResult []map[string]any, err error) {
	if !safeIdentifierRegex.MatchString(schema) || !safeIdentifierRegex.MatchString(tableName) {
		return nil, fmt.Errorf("nombre de esquema o tabla no válido: contiene caracteres no permitidos")
	}

	query := fmt.Sprintf(getDataParameterizedTables, schema, tableName)

	// 3. Ejecutar la consulta
	rows, err := s.db.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}
	defer rows.Close() // ¡Muy importante! Asegura que la conexión se libere.

	// 4. Obtener los nombres de las columnas del resultado
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error al obtener los nombres de las columnas: %w", err)
	}

	var results []map[string]any

	// Iterar sobre cada fila del resultado
	for rows.Next() {
		values := make([]any, len(columns))
		scanArgs := make([]any, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Escanea la fila actual en los punteros
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}

		// Crea un mapa para la fila actual.
		rowMap := make(map[string]any)
		for i, colName := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		// Agrega el mapa de la fila al slice de resultados.
		results = append(results, rowMap)
	}

	// Verifica si hubo algún error durante la iteración de las filas.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de filas: %w", err)
	}

	return results, nil
}

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
