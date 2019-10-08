package accounts

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thiepwong/microservices/common"
)

// AccountService ...interface ...
type AccountService interface {
	Register(*RegisterModel) (interface{}, error)
	Update(profile string) (interface{}, error)
	Profile(id string, token string) (string, error)
}

type accountServiceImp struct {
	repo AccountRepository
	conf *common.Config
}

// NewAccountService ...
func NewAccountService(repo AccountRepository, conf *common.Config) AccountService {
	return &accountServiceImp{repo: repo, conf: conf}
}

func (s *accountServiceImp) Register(register *RegisterModel) (interface{}, error) {

	_code := uuid.Must(uuid.NewV4())
	register.ID = strings.TrimSpace(register.Username)
	register.Username = strings.TrimSpace(register.Username)
	register.VerifyCode = _code.String()
	register.RegisteredDate = time.Now().Unix()
	register.Profile.FullName = fmt.Sprintf("%s %s", register.Profile.LastName, register.Profile.FirstName)
	_pwd, err := common.Hash(register.Password, common.Salt)
	if err != nil {
		return nil, err
	}
	register.Password = _pwd
	_username, err := common.ValidateUsername(register.Username)
	switch _username {
	case 0:
		return nil, err
	case 1:
		// Is email type username
		register.Profile.Email = register.Username
		break
	case 2:
		// Is mobile type username
		register.Profile.Mobile = register.Username
		break
	}

	if err != nil {
		return nil, err
	}
	return s.repo.Register(register)
}

func (s *accountServiceImp) Update(profile string) (interface{}, error) {
	fmt.Println(profile)
	return nil, nil
}

func (s *accountServiceImp) Profile(id string, token string) (string, error) {

	//func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error

	return "Da lay duoc profile tu ID: " + id, nil
}
