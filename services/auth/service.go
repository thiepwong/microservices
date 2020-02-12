package auth

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/thiepwong/microservices/common"
)

const (
	node = "node001"
)

//AuthService interface
type AuthService interface {
	SignIn(*SignInModel) (interface{}, error)
	Verify(activate *ActivateModel) (interface{}, error)
	UpdateContact(*UpdateContact) (interface{}, error)
	ChangePassword(*ChangePasswordModel) (bool, error)
	ConfirmVerify(*VerifyContact) (bool, error)
	SignInViaRefreshToken(string) (interface{}, error)
	ForgotPassword(string, string, string) (bool, error)
}

type authServiceImp struct {
	repo AuthRepository
	conf *common.Config
}

// NewAuthService method
func NewAuthService(repo AuthRepository, conf *common.Config) AuthService {
	return &authServiceImp{repo: repo, conf: conf}
}

func (s *authServiceImp) SignIn(signin *SignInModel) (interface{}, error) {
	acc, prof, err := s.repo.SignIn(signin)
	if err != nil {
		return nil, err
	}

	valid := common.PasswordCompare(signin.Password, acc.Password, common.Salt)
	if valid == false {
		return nil, errors.New("Password is invalid! Please try again!")
	}

	_data, err := json.Marshal(acc)
	err = json.Unmarshal(_data, &acc)
	_iat := time.Now().Unix()
	_exp := _iat

	if signin.Expired > 0 {
		_exp += int64(signin.Expired)
	} else {
		_exp += int64(60 * 60 * 24)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"usr": prof.FullName,
		"iss": node,
		"act": acc.Username,
		"sid": acc.ProfileID,
		"jit": strings.Replace(uuid.Must(uuid.NewV4(), errors.New("error")).String(), "-", "", -1),
		"iat": _iat,
		"exp": _exp,
		"sys": signin.System,
	})

	// fmt.Print(prof)

	rsa, err := common.ReadPrivateKey("./1011.perm")

	tokenString, err := token.SignedString(rsa)
	if err != nil {
		return nil, err
	}

	_refTkn := RefreshToken{
		UserName: acc.Username,
		Profile:  prof,
	}

	_refreshToken, err := createRefreshToken(s, &_refTkn)
	if err != nil {
		return nil, err
	}

	response := &SignInResponse{
		SmartID:      acc.ProfileID,
		Token:        tokenString,
		RefreshToken: _refreshToken,
	}
	return response, nil
}

func (s *authServiceImp) SignInViaRefreshToken(refreshToken string) (interface{}, error) {

	_rft, err := s.repo.ReadRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	var refTkn RefreshToken
	err = json.Unmarshal([]byte(_rft), &refTkn)
	if err != nil {
		return nil, err
	}

	_iat := time.Now().Unix()
	_exp := _iat

	_exp += int64(60 * 60 * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"usr": refTkn.Profile.FullName,
		"iss": node,
		"act": refTkn.UserName,
		"sid": refTkn.Profile.ID,
		"jit": strings.Replace(uuid.Must(uuid.NewV4(), errors.New("error")).String(), "-", "", -1),
		"iat": _iat,
		"exp": int64(60 * 60 * 24),
		"sys": "ref-tk",
	})

	// fmt.Print(prof)

	rsa, err := common.ReadPrivateKey("./1011.perm")

	tokenString, err := token.SignedString(rsa)
	if err != nil {
		return nil, err
	}

	_refreshToken, err := createRefreshToken(s, &refTkn)
	if err != nil {
		return nil, err
	}

	response := &SignInResponse{
		SmartID:      refTkn.Profile.ID,
		Token:        tokenString,
		RefreshToken: _refreshToken,
	}
	return response, nil

}

func (s *authServiceImp) Verify(activate *ActivateModel) (interface{}, error) {

	_userType, err := common.ValidateUsername(activate.Username)
	_register := &RegisterModel{}
	if err != nil {
		return nil, err
	}
	switch _userType {
	case 1:
		_register, err = s.repo.VerifyByEmail(activate.Username, activate.ActivateCode)
		break

	case 2:
		_register, err = s.repo.VerifyBySms(activate.Username, activate.ActivateCode)
		break

	default:
		return nil, errors.New("Username is not valid!")
	}

	// If register not found in system return empty
	if _register == nil {
		return _register, err
	}

	smartID, err := common.GenerateSmartID(8, 1, 16)

	if err != nil {
		return nil, errors.New("Cannot create a Smart ID, please try again!")
	}

	res, err := s.repo.CreateID(_register, smartID)

	return res, err
}

func (s *authServiceImp) UpdateContact(contact *UpdateContact) (interface{}, error) {
	_contact, err := common.ValidateUsername(contact.Contact)
	if err != nil {
		return nil, err
	}
	switch _contact {
	case 1:
		// Update Email to existed profile
		// Load mailpools
		_mailPool, err := s.repo.ReadMailPool(contact.Contact)
		if err != nil {
			return nil, err
		}

		if _mailPool.Code != contact.Code {
			return nil, errors.New("Code not match, please try again later")
		}

		if _mailPool.Used == true {
			// Combine 2 user account with new Smart ID
			return s.repo.UpdateProfileContactWithCombineUser(_mailPool.Email, _mailPool.SID, 1)

		} else {
			// Update mail to profile contact email without combine any user account
			return s.repo.UpdateProfileContact(_mailPool.Email, _mailPool.SID, 1)

		}

	case 2:
		// Update Mobile to existed profile
		// Load mobilepools
		_mobilePool, err := s.repo.ReadMobilePool(contact.Contact, contact.Code)
		if err != nil {
			return nil, err
		}

		if _mobilePool == nil {
			return nil, errors.New("This update confirm is invalid")
		}

		if _mobilePool.Used == true {
			// Update mobile number to profile with combine user account
			return s.repo.UpdateProfileContactWithCombineUser(_mobilePool.Mobile, _mobilePool.SID, 2)
		} else {
			// Update mobile number to profile without any combine

			return s.repo.UpdateProfileContact(_mobilePool.Mobile, _mobilePool.SID, 2)
		}
	default:

		break
	}

	return nil, nil
}

func (s *authServiceImp) ChangePassword(pwd *ChangePasswordModel) (bool, error) {
	_usr, err := s.repo.ReadPassword(pwd.Username)
	if err != nil {
		return false, err
	}

	valid := common.PasswordCompare(pwd.OldPassword, _usr.Password, common.Salt)
	if valid == false {
		return false, errors.New("Old password not match! Please try again!")
	}

	_pwd, err := common.Hash(pwd.NewPassword, common.Salt)
	if err != nil {
		return false, err
	}
	return s.repo.WritePassword(pwd.Username, _pwd)
}

func (s *authServiceImp) ConfirmVerify(cont *VerifyContact) (bool, error) {
	return false, nil
}

func (s *authServiceImp) ForgotPassword(username string, password string, code string) (bool, error) {

	res, err := s.repo.ReadOTP(code)
	if err != nil {
		return false, err
	}
	var _otp = &OtpModel{}
	err = json.Unmarshal([]byte(res), _otp)
	if err != nil {
		return false, err
	}
	_pwd, err := common.Hash(password, common.Salt)
	if err != nil {
		return false, err
	}

	return s.repo.WritePassword(username, _pwd)
}

func createRefreshToken(s *authServiceImp, _rf *RefreshToken) (string, error) {
	_code := strings.Replace(uuid.Must(uuid.NewV4(), errors.New("error")).String(), "-", "", -1)
	_json, err := json.Marshal(_rf)
	if err != nil {
		return "", err
	}

	_, err = s.repo.CreateRefreshToken(_code, _json, 60*60*24*180)
	if err != nil {
		return "", errors.New("Create refresh token is failed")
	}
	return _code, nil
}
