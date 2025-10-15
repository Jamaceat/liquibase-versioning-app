package dto

type (
	TableName struct {
		TableName  string `json:"table_name"`
		SchemaName string `json:"schema_name"`
	}

	TableDetailed struct {
		TableName  string `json:"table_name"`
		SchemaName string `json:"schema_name"`
		Columns    []Column
	}

	Column struct {
		ColumnName  string `json:"column_name"`
		NoAliasType string `json:"tipo_de_dato_completo"`
		AliasType   string `json:"alias_del_tipo"`
	}
)
