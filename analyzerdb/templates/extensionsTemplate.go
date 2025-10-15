package templates

const (
	ExtensionTemplate = `

	-- Extension %s version %s
	CREATE EXTENSION IF NOT EXISTS %s WITH VERSION '%s';
	
	`
)
