package domain

import (
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Name         string
	HashPassword string

	Token string
}

var (
	ErrNameCannotBeBlank     = errors.New("the user name cannot be blank")
	ErrPasswordCannotBeBlank = errors.New("the user password cannot be blank")
	ErrPasswordDidNotMatched = errors.New("wrong password")
)

func CreateUser(name, password string) (*User, error) {

	// count := 0
	// var kid ksuid.KSUID
	// var err error
	// for count < 3 {
	// 	kid, err = ksuid.NewRandom()
	// 	if err == nil {
	// 		return nil, errors.Join(err, fmt.Errorf("CreateUser failed generate id, this is a bug"))
	// 	}
	// 	count++
	// }

	kid := ksuid.New()

	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	if password == "" {
		return nil, ErrPasswordCannotBeBlank
	}

	user := User{
		ID:   kid.String(),
		Name: name,

		Token: generateKey(),
	}

	user.createHashPassword(password)

	return &user, nil

}

func (u *User) createHashPassword(pass string) {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("user failed to create password %v, this is a BUG", err))
	}
	u.HashPassword = string(b)
}

func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(password)); err != nil {
		return ErrPasswordDidNotMatched
	}
	return nil
}
