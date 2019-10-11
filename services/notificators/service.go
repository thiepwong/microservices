package notificators

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
	_url "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thiepwong/microservices/common"
)

type NotificatorService interface {
	SendEmail(*MailModel) (bool, error)
	SendSMS(*SmsModel) (interface{}, error)
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

	var r *common.Request
	_host := fmt.Sprintf("%s:%d", s.conf.Option.EmailSender.Server, s.conf.Option.EmailSender.Port)
	auth := smtp.PlainAuth("", s.conf.Option.EmailSender.Email, s.conf.Option.EmailSender.Password, s.conf.Option.EmailSender.Server)
	r = common.NewRequest([]string{mail.Mail}, mail.Subject, mail.Content, _host, s.conf.Option.EmailSender.Email, auth)

	switch mail.Type {

	case 3:
		res, err := s.repo.ReadMailPool(mail.Mail)
		if err != nil {
			return false, errors.New("This email is not found in system")
		}

		templateData := struct {
			FullName string
			SID      uint64
			Email    string
			URL      string
		}{
			FullName: res.FullName,
			SID:      res.SID,
			Email:    res.Email,
			URL:      fmt.Sprintf("%s?contact=%s&code=%s", s.conf.Option.UpdateContactURL, res.Email, res.Code),
		}
		err = r.ParseTemplate(fmt.Sprintf("./templates/%s/verify_email.html", mail.Lang), templateData)
		if err != nil {
			return false, err
		}

		break
	case 5:
		res, err := s.repo.ReadMailActivatedCode(mail.Mail)
		if err != nil {
			return false, errors.New("This username is not found in system")
		}

		if res.VerifyCode == "" {
			return false, errors.New("This account is activated before!")
		}

		templateData := struct {
			FullName string
			Username string
			URL      string
		}{
			FullName: res.Profile.FullName,
			Username: res.Username,
			URL:      fmt.Sprintf("%s?username=%s&code=%s", s.conf.Option.ActivateURL, res.Username, res.VerifyCode),
		}
		err = r.ParseTemplate(fmt.Sprintf("./templates/%s/activate_email.html", mail.Lang), templateData)
		if err != nil {
			return false, err
		}
		break
	default:
		break
	}

	ok, _ := r.SendEmail()

	return ok, nil

	//	return sendMail(mail.Mail, mail.Subject, mail.Content, s.conf.Option.EmailSender)
}

func (s *notificatorServiceImpl) SendSMS(sms *SmsModel) (interface{}, error) {

	_irisToken := s.repo.ReadIrisToken(s.conf.SmsIris.Brandname)
	if _irisToken == "" {
		// Sign In and save to Redis
		_iris, err := smsSignin(&s.conf.SmsIris)
		if err != nil {
			return nil, err
		}

		_irisToken = _iris.AccessToken
		s.repo.WriteIrisToken(s.conf.SmsIris.Brandname, _irisToken, _iris.Expired)
	}

	// Check Send SMS type
	// 1 => Created OTP then send it with default template
	// 2 => Created OTP then send it with custom template
	// 3 => Send a normal sms without any OTP for registered member
	// 4 => Send to any one with any content
	var _otp *OtpModel
	if sms.Lang == "" {
		sms.Lang = "vi"
	}
	switch sms.Type {
	case 1:
		_, err := s.repo.ReadRegisterByUser(sms.Mobile)
		if err != nil {
			return nil, errors.New("This mobile number is not registered for SmartID")
		}
		if sms.TTL == 0 {
			sms.TTL = 120
		}
		_template, err := common.ReadTemplate(fmt.Sprintf("./templates/%s/activate_otp.msg", sms.Lang))
		if err != nil {
			return nil, errors.New("Sms template not found! Please use custom template to send")
		}

		_otp = generateOTP(sms.Mobile, 6, sms.TTL)
		_sms := fmt.Sprintf(_template, _otp.ID, _otp.TTL)
		_json, err := json.Marshal(_otp)
		if err != nil {
			return false, err
		}
		s.repo.SaveOTP(_otp.ID, _json, _otp.TTL)

		res, err := smsSender(sms.Mobile, _sms, _irisToken, &s.conf.SmsIris)
		if err != nil {
			return nil, err
		}

		return res, nil

	case 2:
		if sms.Content == "" {
			return nil, errors.New("Content template is required")
		}
		_, err := s.repo.ReadRegisterByUser(sms.Mobile)
		if err != nil {
			return nil, errors.New("This mobile number is not registered for SmartID")
		}
		if sms.TTL == 0 {
			sms.TTL = 120
		}
		_otp = generateOTP(sms.Mobile, 6, sms.TTL)
		_sms := fmt.Sprintf(sms.Content, _otp.ID, _otp.TTL)
		_json, err := json.Marshal(_otp)
		if err != nil {
			return false, err
		}
		s.repo.SaveOTP(_otp.ID, _json, _otp.TTL)

		res, err := smsSender(sms.Mobile, _sms, _irisToken, &s.conf.SmsIris)
		if err != nil {
			return nil, err
		}

		return res, nil

	case 3:
		if sms.Content == "" {
			return nil, errors.New("Content is required!")
		}

		_, err := s.repo.ReadRegisterByUser(sms.Mobile)
		if err != nil {
			return nil, errors.New("This mobile number is not registered for SmartID")
		}

		res, err := smsSender(sms.Mobile, sms.Content, _irisToken, &s.conf.SmsIris)
		if err != nil {
			return nil, err
		}
		return res, nil

	case 4:
		if sms.Content == "" {
			return nil, errors.New("Content is required!")
		}
		res, err := smsSender(sms.Mobile, sms.Content, _irisToken, &s.conf.SmsIris)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, nil
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

// Deprecated method by a html emailer
func sendMail(email string, subject string, content string, mailConf *common.MailSender) (bool, error) {
	err := smtp.SendMail(fmt.Sprintf("%s:%d", mailConf.Server, mailConf.Port),
		smtp.PlainAuth("", mailConf.Email, mailConf.Password, mailConf.Server),
		mailConf.Email, []string{email}, []byte(content))

	if err != nil {
		return false, err
	}

	return true, nil
}

func smsSignin(cfg *common.IRIS) (iris *IrisSignInResponse, err error) {

	form := _url.Values{}
	form.Add("grant_type", "password")

	authen := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.Username, cfg.Password)))
	url := fmt.Sprintf("%s%s", cfg.Host, cfg.SignInRoute)
	req, e := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if e != nil {
		return nil, errors.New("Cannot connect SMS server")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authen))

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Not response from SMS server")
	}

	_body, err := ioutil.ReadAll(response.Body)
	var _res *IrisSignInResponse = &IrisSignInResponse{}

	err = json.Unmarshal(_body, _res)

	if err != nil {
		return nil, err
	}

	return _res, nil
}

func smsSender(mobile string, content string, token string, cfg *common.IRIS) (interface{}, error) {
	var _msg SenderModel = SenderModel{
		Brandname: cfg.Brandname,
		SendingList: []Message{Message{
			SmsId:       fmt.Sprintf("%s-%d", mobile, time.Now().Unix()),
			PhoneNumber: mobile,
			Content:     content,
			ContentType: "30",
		}},
	}

	_msgByte, err := json.Marshal(_msg)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", cfg.Host, cfg.SendSmsRoute)
	req, e := http.NewRequest("POST", url, bytes.NewBuffer(_msgByte))
	if e != nil {
		return nil, errors.New("Cannot connect SMS server")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer  %s", token))

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Not response from SMS server")
	}

	_body, err := ioutil.ReadAll(response.Body)
	var _res *IrisSentResponse = &IrisSentResponse{}

	err = json.Unmarshal(_body, _res)

	if err != nil {
		return nil, err
	}
	if _res.Code != "201" {
		return nil, errors.New(_res.Message)
	}
	return _res, nil

}
