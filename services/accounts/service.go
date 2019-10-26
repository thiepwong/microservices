package accounts

import (
	"errors"
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
	Profile(id uint64, token string) (interface{}, error)
	UpdateEmail(prof *AuthUpdate) (bool, error)
	UpdateMobile(prof *AuthUpdate) (bool, error)
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
	register.ID = strings.TrimSpace(register.Username)
	register.Username = strings.TrimSpace(register.Username)
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
		_code := uuid.Must(uuid.NewV4())
		register.VerifyCode = _code.String()
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

func (s *accountServiceImp) Profile(sid uint64, token string) (interface{}, error) {

	return s.repo.GetProfileById(sid)
}

func (s *accountServiceImp) UpdateEmail(prof *AuthUpdate) (bool, error) {
	// Check the profile id is valid?
	_p, err := s.repo.GetProfileById(prof.SID)
	if err != nil {
		return false, err
	}

	if _p.Email == prof.Email {
		return false, errors.New("This email was added before!")
	}
	// Check email was registered before or not!
	_user, err := s.repo.GetUserById(prof.Email)
	if err != nil {
		return false, err
	}

	_code := uuid.Must(uuid.NewV4())
	var _m EmailProfileModel
	if _user.ID == "" {
		// Cho phep update du lieu vao profile
		_m = EmailProfileModel{
			ID:       prof.Email,
			SID:      _p.ID,
			Code:     _code.String(),
			FullName: _p.FullName,
			Username: prof.Username,
			Email:    prof.Email,
			Used:     false,
		}

	} else {
		_m = EmailProfileModel{
			ID:       prof.Email,
			SID:      _p.ID,
			Code:     _code.String(),
			Username: prof.Username,
			Email:    prof.Email,
			Used:     true,
		}
	}

	return s.repo.CreateEmailPool(&_m)
	// Create new user with same profile id
}
func (s *accountServiceImp) UpdateMobile(prof *AuthUpdate) (bool, error) {

	// Check the profile id is valid?
	_p, err := s.repo.GetProfileById(prof.SID)
	if err != nil {
		return false, err
	}

	if _p.Mobile == prof.Mobile {
		return false, errors.New("This mobile was added before!")
	}

	// Check email was registered before or not!
	_user, err := s.repo.GetUserById(prof.Mobile)
	if err != nil {
		return false, err
	}

	var _m MobileProfileModel
	if _user.ID == "" {
		// Cho phep update du lieu vao profile
		_m = MobileProfileModel{
			ID:       prof.Mobile,
			SID:      _p.ID,
			Code:     common.GenerateOTP(prof.Mobile, 6, 120).ID,
			FullName: _p.FullName,
			Username: prof.Username,
			Mobile:   prof.Mobile,
			Used:     false,
		}

	} else {
		_m = MobileProfileModel{
			ID:       prof.Mobile,
			SID:      _p.ID,
			Code:     common.GenerateOTP(prof.Mobile, 6, 120).ID,
			FullName: _p.FullName,
			Username: prof.Username,
			Mobile:   prof.Mobile,
			Used:     true,
		}
	}

	return s.repo.CreateMobilePool(&_m)
}
