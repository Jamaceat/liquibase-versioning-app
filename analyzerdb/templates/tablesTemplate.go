package templates

const (
	TableTemplate = `

	CREATE TABLE IF NOT EXISTS "%s".%s ( 
	%s
	);

	`

	ColumnTableTemplate = `%s %s`
)
