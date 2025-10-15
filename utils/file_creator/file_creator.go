package filecreator

import (
	"fmt"
	"io"
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

	for _, value := range typeEntity.TypeEntityColumns {

		typeStructureSetted += fmt.Sprintf(templates.TypeStructureTemplate, value.ColumnName, value.AliasType)

	}

	formatted = fmt.Sprintf(templates.TypeTemplate,
		typeEntity.TypeName.TypName,
		typeEntity.SchemaName,
		typeEntity.TypeName.TypName,
		typeStructureSetted,
	)

	return
}
