package auth

import (
	"github.com/thiepwong/microservices/common"
)

// AccountService ...interface ...
type AuthService interface {
	SignIn(*SignInModel) (interface{}, error)
}

type authServiceImp struct {
	repo AuthRepository
	conf *common.Config
}

// NewAccountService ...
func NewAuthService(repo AuthRepository, conf *common.Config) AuthService {
	return &authServiceImp{repo: repo, conf: conf}
}

func (s *authServiceImp) SignIn(signin *SignInModel) (interface{}, error) {

	return s.repo.SignIn(signin)
}
