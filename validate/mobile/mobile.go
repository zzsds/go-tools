package mobile

import "regexp"

// ValidateMobile ...
func Validate(mobile string) bool {
	ok, _ := regexp.MatchString(`^((\+[0-9]\d{10,12})|1[1-9]\d{9})$`, mobile)
	return ok
}
