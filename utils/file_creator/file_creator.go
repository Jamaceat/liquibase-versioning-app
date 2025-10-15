package filecreator

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/dto"
	"github.com/Jamaceat/liquibase-versioning-app/analyzerdb/templates"
)

func GenerateSQLFile(writer io.Writer, script string) error {
	_, err := fmt.Fprintf(writer, "-- Auto-generated SQL Backup\n-- Date: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(writer, script)
	if err != nil {
		return err
	}

	return nil
}

func FormatType(typeEntity dto.TypesEntityComplete) (formatted string) {

	var typeStructureSetted string

	var columns []string
	for _, value := range typeEntity.TypeEntityColumns {

		columns = append(columns, fmt.Sprintf(templates.TypeStructureTemplate, value.ColumnName, value.AliasType))

	}

	typeStructureSetted = strings.Join(columns, ",\n\t\t")

	formatted = fmt.Sprintf(templates.TypeTemplate,
		typeEntity.TypeName.TypName,
		typeEntity.SchemaName,
		typeEntity.TypeName.TypName,
		typeStructureSetted,
	)

	return
}

func FormatEnum(enumEntityName string, enumEntityValues []string, schema string) string {

	var valuesMap []string = make([]string, 0)

	for _, value := range enumEntityValues {
		valuesMap = append(valuesMap, fmt.Sprintf("'%s'", value))
	}

	enumValues := strings.Join(valuesMap, ",")

	structuredParameters := fmt.Sprintf(templates.EnumStructureTemplate, enumValues)

	return fmt.Sprintf(templates.EnumTemplate, enumEntityName, schema, enumEntityName, structuredParameters)

}
