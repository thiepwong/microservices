package images

import (
	"errors"
	"fmt"

	"github.com/thiepwong/microservices/common"
)

// AccountService ...interface ...
type ImageService interface {
	Upload(*Image, uint64) (interface{}, error)
	List(uint64) (interface{}, error)
}

type imageServiceImp struct {
	conf *common.Config
}

// NewAccountService ...
func NewImageService(conf *common.Config) ImageService {
	return &imageServiceImp{conf: conf}
}

func (s *imageServiceImp) Upload(img *Image, sid uint64) (interface{}, error) {
	if img.Data == "" {
		return nil, errors.New("Data is empty! Please upload an image")
	}
	name, err := common.FtpWriteImage(s.conf, fmt.Sprintf("%d", sid), img.Data)
	if err != nil {
		return nil, err
	}

	return name, nil

}

func (s *imageServiceImp) List(sid uint64) (interface{}, error) {
	list, err := common.FtpListAll(s.conf, fmt.Sprintf("%d", sid))
	return list, err
}
