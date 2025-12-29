package models

import "regexp"

type CategoryRule struct {
	Pattern  *regexp.Regexp
	Category string
}
