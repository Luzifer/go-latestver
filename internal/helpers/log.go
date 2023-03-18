// Package helpers contains helper functions to avoid code duplication
package helpers

import "github.com/sirupsen/logrus"

// LogIfErr yields a logrus error log line when the given error is
// not nil
func LogIfErr(err error, msgTpl string, params ...any) {
	if err == nil {
		return
	}

	logrus.WithError(err).Errorf(msgTpl, params...)
}
