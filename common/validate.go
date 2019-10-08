package common

import (
	"errors"
	"regexp"
)

func ValidateUsername(username string) (int, error) {
	if username == "" {
		return 0, errors.New("username is required!")
	}
	emailReg := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isEmail := emailReg.MatchString(username)
	if isEmail == true {
		return 1, nil
	} else {

		mobileReg := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
		isMobile := mobileReg.MatchString(username)
		if isMobile == true {
			return 2, nil
		} else {
			return 0, errors.New("Username is invalid!")
		}

	}
}
