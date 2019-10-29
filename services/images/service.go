package images

import (
	"errors"

	"github.com/thiepwong/microservices/common"
)

// AccountService ...interface ...
type AccountService interface {
	Upload(*Image) (interface{}, error)
}

type accountServiceImp struct {
	repo AccountRepository
	conf *common.Config
}

// NewAccountService ...
func NewAccountService(repo AccountRepository, conf *common.Config) AccountService {
	return &accountServiceImp{repo: repo, conf: conf}
}

func (s *accountServiceImp) Upload(img *Image) (interface{}, error) {
	if img.Data == "" {
		return nil, errors.New("Data is empty! Please upload an image")
	}

	return img, nil

}
