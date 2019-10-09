package notificators

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/thiepwong/microservices/common"
)

type NotificatorService interface {
	SendEmail(*MailModel) (bool, error)
	SendSMS(mobile string, content string, ttl time.Duration) (interface{}, error)
	SendFirebase(channel string, title string, content string) (interface{}, error)
}

type notificatorServiceImpl struct {
	repo NotificatorRepository
	conf *common.Config
}

func NewNotificatorService(repo NotificatorRepository, conf *common.Config) NotificatorService {
	return &notificatorServiceImpl{
		repo: repo,
		conf: conf,
	}
}

func (s *notificatorServiceImpl) SendEmail(mail *MailModel) (bool, error) {

	switch mail.Type {
	case 5:
		res, err := s.repo.ReadMailActivatedCode(mail.Mail)
		if err != nil {
			return false, errors.New("This username is not found in system")
		}

		if res.VerifyCode == "" {
			return false, errors.New("This account is activated before!")
		}

		mail.Content = fmt.Sprintf("Day la ma so kich hoat cua he thong smart id %s cho tai khoan %s", res.VerifyCode, res.Username)
		break
	default:
		break
	}

	return sendMail(mail.Mail, mail.Subject, mail.Content, s.conf.Option.EmailSender)
}

func (s *notificatorServiceImpl) SendSMS(mobile string, content string, ttl time.Duration) (interface{}, error) {
	_otp := generateOTP(mobile, 6, ttl)

	_json, err := json.Marshal(_otp)
	if err != nil {
		return false, err
	}
	// using sms microservice to send this otp
	res, err := s.repo.SaveOTP(_otp.ID, string(_json), ttl)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (s *notificatorServiceImpl) SendFirebase(channel string, title string, content string) (interface{}, error) {
	return nil, nil
}

func generateOTP(mobile string, size int, ttl time.Duration) *OtpModel {
	var code string
	for i := 0; i < size; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	_otp := &OtpModel{ID: code, Mobile: mobile, TTL: ttl}
	return _otp
}

func sendMail(email string, subject string, content string, mailConf *common.MailSender) (bool, error) {
	err := smtp.SendMail(fmt.Sprintf("%s:%d", mailConf.Server, mailConf.Port),
		smtp.PlainAuth("", mailConf.Email, mailConf.Password, mailConf.Server),
		mailConf.Email, []string{email}, []byte(content))

	if err != nil {
		return false, err
	}

	return true, nil
}
