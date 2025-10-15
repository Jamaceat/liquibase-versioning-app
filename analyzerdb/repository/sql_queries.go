package repository

const (
	getExtensions = `
		SELECT pe.extname, pe.extversion FROM pg_catalog.pg_extension pe 
	`

	getTypes = `
	SELECT * FROM pg_catalog.pg_type pt WHERE pt.typowner != 10 --typowner 10 es por defecto
	AND pt.typname ~ '^[a-zA-Z].+' 
	AND pt.typnamespace  IN (SELECT pn."oid"  FROM pg_catalog.pg_namespace pn WHERE pn.nspname ='%s') --schema
	AND pt.typname  IN (SELECT pc.relname FROM pg_catalog.pg_class pc WHERE pc.relkind ='c')
	`

	getTypesDetailed = `
		SELECT
		a.attname AS column_name,
		format_type(a.atttypid, a.atttypmod) AS no_alias,
		pt_attr.typname AS alias
	FROM
		pg_catalog.pg_type pt_main
	JOIN
		pg_catalog.pg_namespace n ON n.oid = pt_main.typnamespace
	JOIN
		pg_catalog.pg_class pc ON pc.oid = pt_main.typrelid
	JOIN
		pg_catalog.pg_attribute a ON a.attrelid = pc.oid
	JOIN
		pg_catalog.pg_type pt_attr ON pt_attr.oid = a.atttypid
	WHERE
		pt_main.typname = '%s' -- nombre del type  
		AND n.nspname = '%s'  --schema   
		AND a.attnum > 0  -- ignorar por defecto los del sistema
	ORDER BY
		a.attnum
	`

	getEnumDetailed = `
		SELECT
		n.nspname AS schema_name,
		t.typname AS enum_name,
		e.enumlabel AS enum_value
		FROM pg_catalog.pg_type t
		JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
		JOIN pg_catalog.pg_enum e ON t.oid = e.enumtypid
		WHERE
		n.nspname = '%s'  -- Filtra por el esquema
		AND t.typtype = 'e'     -- Filtra solo los tipos que son ENUM
		ORDER BY 
		enum_name, e.enumsortorder
	`

	getSequencesDetailed = `
		SELECT
		n.nspname AS nombre_esquema,
		cl_tabla.relname AS nombre_tabla,
		a.attname AS nombre_columna,
		cl_seq.relname AS nombre_secuencia,
		format_type(a.atttypid, a.atttypmod) AS tipo_de_dato_completo,
		t.typname AS alias_del_tipo
	FROM
		pg_depend AS d
	JOIN
		pg_class AS cl_tabla ON cl_tabla.oid = d.refobjid
	JOIN
		pg_attribute AS a ON a.attrelid = cl_tabla.oid AND a.attnum = d.refobjsubid
	JOIN
		pg_class AS cl_seq ON cl_seq.oid = d.objid
	JOIN
		pg_catalog.pg_type AS t ON t.oid = a.atttypid
	JOIN
		pg_catalog.pg_namespace AS n ON n.oid = cl_tabla.relnamespace
	WHERE
		cl_tabla.relkind = 'r' -- El objeto referenciado es una tabla
		AND cl_seq.relkind = 'S' -- El objeto dependiente es una secuencia ('S')
		AND n.nspname = '%s' --schema
	`

	getTables = `
			SELECT * FROM pg_catalog.pg_tables pt WHERE pt.schemaname ='%s' --schema
	`

	getTableDetail = `
	SELECT
		t.relname AS nombre_tabla,
		a.attname AS nombre_columna,
		format_type(a.atttypid, a.atttypmod) AS tipo_de_dato_completo,
		pt.typname AS alias_del_tipo
	FROM
		pg_catalog.pg_class AS t
	JOIN
		pg_catalog.pg_attribute AS a ON a.attrelid = t.oid
	JOIN
		pg_catalog.pg_type AS pt ON pt.oid = a.atttypid
	JOIN
		pg_catalog.pg_namespace AS n ON n.oid = t.relnamespace
	WHERE
		t.relkind = 'r' -- Asegura que solo sean tablas
		AND n.nspname = 'public' -- Filtra por el esquema public
		AND a.attnum > 0 -- Excluye columnas del sistema
		AND NOT a.attisdropped -- Excluye columnas eliminadas
		AND t.relname = '%s' -- nombre de la tabla
	ORDER BY
		t.relname,
		a.attnum
	`

	getTriggers = `
	SELECT
	pn.nspname AS table_schema,
	pc.relname AS table_name,
	pt.tgname AS trigger_name,
	pg_get_triggerdef(pt.oid) AS trigger_definition
	FROM pg_catalog.pg_trigger pt
	JOIN pg_catalog.pg_class pc
	ON pt.tgrelid = pc.oid
	JOIN pg_catalog.pg_namespace pn
	ON pc.relnamespace = pn.oid
	WHERE
	pn.nspname = '%s' -- schema
	`
)
