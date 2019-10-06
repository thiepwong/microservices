package accounts

import "fmt"

// AccountService ...interface ...
type AccountService interface {
	Register(username string, password string, profile string) (interface{}, error)
	Update(profile string) (interface{}, error)
	Profile(id string) (string, error)
}

type accountServiceImp struct {
	db string
}

// NewAccountService ...
func NewAccountService(profile string) AccountService {
	fmt.Println(profile)

	return &accountServiceImp{db: profile}
}

func (s *accountServiceImp) Register(username string, password, profile string) (interface{}, error) {
	fmt.Println(s.db)
	return nil, nil
}

func (s *accountServiceImp) Update(profile string) (interface{}, error) {
	fmt.Println(profile)
	return nil, nil
}

func (s *accountServiceImp) Profile(id string) (string, error) {
	return "Da lay duoc profile tu ID: " + id, nil
}
