package server

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func errorJSON(err error) []byte {
	log.WithFields(logrus.Fields{"error": err}).Error()
	return []byte(fmt.Sprintf(`{"result":"","error":%q}`, err))
}
