package templates

const (
	ExtensionTemplate = `

	-- Extension %s version %s
	CREATE EXTENSION IF NOT EXISTS %s WITH VERSION '%s';
	
	`
	//%s
	//	1. TypeName
	//	2. SchemaName
	//	3. TypeName
	//	4. EnumStructure>>
	EnumTemplate = `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN  -- enum name
			CREATE TYPE %s."%s" AS ENUM %s;  -- schema , typename,
		END IF;
	END$$;
	`

	//%s
	//	1. 'enum1','enum2','enum3'
	EnumStructureTemplate = `(%s)`

	// %s
	//
	// 	1. typename
	// 	2. schemaName
	// 	3. typename
	// 	4. TypeStructureTemplate >>
	TypeTemplate = `
	
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
				CREATE TYPE %s.%s AS (
					%s
				);
			END IF;
		END$$;

	`
	// %s
	//	ColumnTypeName
	//	DataTypeAlias
	TypeStructureTemplate = `"%s" %s
	`
)
