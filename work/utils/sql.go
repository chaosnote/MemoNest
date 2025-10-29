package utils

import (
	"fmt"
	"regexp"
)

func ParseSQLError(err error, default_message string) error {
	if err == nil {
		return err
	}

	re := regexp.MustCompile(`\s*\[ERR\](.*)`)
	match := re.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		return fmt.Errorf("%s", match[1])
	} else {
		return fmt.Errorf("%s", default_message)
	}
}
