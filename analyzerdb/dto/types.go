package dto

type (
	TypesEntityN struct {
		TypName string `json:"typname"`
	}

	TypesEntityComplete struct {
		TypeName          TypesEntityN
		TypeEntityColumns []TypeEntityColumns
		SchemaName        string
	}

	TypeEntityColumns struct {
		ColumnName string `json:"column_name"`
		AliasType  string `json:"alias"`
	}

	EnumEntity struct {
		SchemaName string `json:"schema_name"`
		EnumName   string `json:"enum_name"`
		EnumValue  string `json:"enum_value"`
	}
)
