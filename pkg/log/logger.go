package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger used in the log package.
var Logger logrus.Logger = *logrus.New()

// SetLevel sets the verbosity of the global logger from a string value such as
// "info". If the value given can not be parsed as a log level, an error is
// returned.
func SetLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("given log level is not a valid level string")
	}
	Logger.SetLevel(lvl)
	return nil
}
