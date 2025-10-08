package utils

import (
	"strings"
)

// sep .
//
//	where<to/from>
//	action
//	who<opt>[uid|*|empty]
func Subject(params ...string) string {
	return strings.Join(params, ".")
}
