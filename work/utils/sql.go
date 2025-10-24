package utils

import (
	"fmt"
	"regexp"
)

func ParseSQLError(err error) error {
	re := regexp.MustCompile(`\s*\[ERR\](.*)`)
	match := re.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		err = fmt.Errorf("%s", match[1])
	}
	return err
}
