package models

import(
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init(){
	logger.Level = logrus.TraceLevel
}
