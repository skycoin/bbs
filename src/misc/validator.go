package misc

import (
	"regexp"
	"errors"
)

func  CheckAlias(s string) (bool,error) {
	if len(s) < 5 || len(s)> 40 {
		return false, errors.New("alias is too short or too long")
	}
	re := regexp.MustCompile("/^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$/")
	if !re.Match([]byte(s)) {
		return false, errors.New("username is invalidat")
	}
	return true,nil
}
