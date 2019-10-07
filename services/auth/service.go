package auth

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	res, err := s.repo.SignIn(signin)
	if err != nil {
		return nil, err
	}

	var _res SignInModel

	_data, err := json.Marshal(res)
	err = json.Unmarshal(_data, &_res)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": _res.Username,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	rsa, err := common.RsaConfigSetup("./1011.perm", "", "./pub.key")

	tokenString, err := token.SignedString(rsa)
	if err != nil {
		return nil, err
	}

	return tokenString, nil
}
