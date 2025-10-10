package filecreator

import (
	"fmt"
	"io"
	"time"
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
