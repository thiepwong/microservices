package accounts

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thiepwong/microservices/common"
)

// AccountService ...interface ...
type AccountService interface {
	Register(*RegisterModel) (interface{}, error)
	Update(profile string) (interface{}, error)
	Profile(id string) (string, error)
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
	_pwd, err := common.Hash(register.Password, 11)
	if err != nil {
		return nil, err
	}
	register.Password = _pwd
	reg := usrValidate(register)
	return s.repo.Register(reg)
}

func (s *accountServiceImp) Update(profile string) (interface{}, error) {
	fmt.Println(profile)
	return nil, nil
}

func (s *accountServiceImp) Profile(id string) (string, error) {
	return "Da lay duoc profile tu ID: " + id, nil
}

func usrValidate(reg *RegisterModel) *RegisterModel {
	res := validateEmail(reg.Username)
	if res == true {
		reg.Profile.Email = reg.Username
	} else {
		reg.Profile.Mobile = reg.Username
	}
	return reg
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
