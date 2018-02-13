package server

import (
	"fmt"
)

func errorJSON(err error) []byte {
	return []byte(fmt.Sprintf(`{"result":"","error":%q}`, err))
}
