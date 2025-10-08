package constant

import (
	"fmt"

	"github.com/Jamaceat/liquibase-versioning-app/utils/path"
)

var ENV_FILE string

func init() {

	ENV_FILE = fmt.Sprintf("%s/.env", path.GetFindProjectRoot())

}
