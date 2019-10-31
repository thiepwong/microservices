package common

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	uuid "github.com/satori/go.uuid"
)

func FtpWriteText(cfg *Config, folder string, data string) (string, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", cfg.Database.FTP.Host, cfg.Database.FTP.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return "", err
	}

	err = c.Login(cfg.Database.FTP.Username, cfg.Database.FTP.Password)
	if err != nil {
		return "", err
	}
	_path := fmt.Sprintf("%s/%s", cfg.Database.FTP.Volume, folder)
	err = c.ChangeDir(_path)
	if err != nil {
		err = c.MakeDir(_path)
		if err != nil {
			return "", err
		}
		c.ChangeDir(_path)
	}

	str := bytes.NewBufferString(data)
	_name := fmt.Sprintf("%s.bin", strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1))
	err = c.Stor(_name, str)
	if err != nil {
		return "", err
	}
	return _name, nil

}

func FtpWriteImage(cfg *Config, folder string, data string) (string, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", cfg.Database.FTP.Host, cfg.Database.FTP.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return "", err
	}

	err = c.Login(cfg.Database.FTP.Username, cfg.Database.FTP.Password)
	if err != nil {
		return "", err
	}
	_path := fmt.Sprintf("%s/%s", cfg.Database.FTP.Volume, folder)
	err = c.ChangeDir(_path)
	if err != nil {
		err = c.MakeDir(_path)
		if err != nil {
			return "", err
		}
		c.ChangeDir(_path)
	}
	_name := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	return saveImage(_name, data, c)

}

func FtpListAll(cfg *Config, folder string) (interface{}, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", cfg.Database.FTP.Host, cfg.Database.FTP.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(cfg.Database.FTP.Username, cfg.Database.FTP.Password)
	if err != nil {
		return nil, err
	}
	_path := fmt.Sprintf("%s/%s", cfg.Database.FTP.Volume, folder)

	list, err := c.List(_path)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func saveImage(name string, data string, c *ftp.ServerConn) (string, error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", errors.New("Invalid image data")
	}
	ImageType := data[11:idx]

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return "", errors.New("Cannot decode base64 data")
	}
	r := bytes.NewReader(unbased)
	switch ImageType {
	case "png":
		_name := fmt.Sprintf("%s.png", name)
		err = c.Stor(_name, r)
		if err != nil {
			return "", err
		}
		return _name, nil
	case "jpeg":
		_name := fmt.Sprintf("%s.jpg", name)
		err = c.Stor(_name, r)
		if err != nil {
			return "", err
		}
		return _name, nil

	case "gif":
		_name := fmt.Sprintf("%s.gif", name)
		err = c.Stor(_name, r)
		if err != nil {
			return "", err
		}
		return _name, nil
	}
	return name, nil
}
