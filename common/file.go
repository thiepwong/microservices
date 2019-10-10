package common

import (
	"errors"
	"io/ioutil"
)

func ReadTemplate(file string) (string, error) {
	_file, err := ioutil.ReadFile(file)
	if err != nil {
		return "", errors.New("File is invalid!")
	}
	return string(_file), nil
}
