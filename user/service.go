package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
}

// private
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	//buat objek dulu
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	// user.Token = helper.TokenString(64);

	if err != nil {
		return user, err
	}
	//diubah ke string karena sebelumnya password bertipe byte
	user.Password = string(password)
	user.Role = "USER"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error){
	email := input.Email
	password := input.Password

	//cari user dengan email yg dimasukan
	user, err:= s.repository.FindByEmail(email)
	if err != nil{
		return user, err
	}

	if user.ID == 0{
		return user, errors.New("No user found with that email")
	}
	//bandingkan hash password dengan password
	//parameter satu ambil data hashed password, parameter kedua password yg dimasukkan
	 err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	 if err != nil{
		return user, err
	 }

	 return user, nil
}

