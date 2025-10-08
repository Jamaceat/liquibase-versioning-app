package path

import (
	"errors"
	"os"
	"path/filepath"
)

func findProjectRoot() (string, error) {
	// Obtiene el directorio de trabajo actual
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Itera hacia arriba desde el directorio actual
	for {
		// Construye la ruta al archivo go.mod en el directorio actual
		goModPath := filepath.Join(wd, "go.mod")

		// Comprueba si el archivo go.mod existe
		if _, err := os.Stat(goModPath); err == nil {
			// ¡Lo encontramos! Devuelve el directorio actual como la raíz del proyecto.
			return wd, nil
		}

		// Sube al directorio padre
		parent := filepath.Dir(wd)

		// Si ya no podemos subir más (llegamos a la raíz del sistema), paramos.
		if parent == wd {
			break
		}
		wd = parent
	}

	return "", errors.New("raíz del proyecto no encontrada (no se encontró go.mod)")
}

func GetFindProjectRoot() (path string) {
	path, err := findProjectRoot()

	if err != nil {
		panic("Not exists go.mod")
	}

	return
}
