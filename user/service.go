package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input EmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserById(id int) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(input FormUpdateUserInput) (User, error)
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

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	//cari user dengan email yg dimasukan
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found with that email")
	}
	//bandingkan hash password dengan password
	//parameter satu ambil data hashed password, parameter kedua password yg dimasukkan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input EmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	//kalo ditemukan email
	if err != nil {
		return false, err
	}
	//kalo emailnya belum ada
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	//dapatkan service berdasarkan ID
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}
	//proses update avatar filena
	user.Avatar = fileLocation 

	//update data ke tabel
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserById(id int) (User, error){
	user, err := s.repository.FindById(id)

	if err != nil{
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found")
	}
	return user, nil
}

func (s *service) GetAllUsers() ([]User, error){
	users, err := s.repository.FindAll()

	if err != nil{
		return users, err
	}

	return users, nil
}

func (s *service) UpdateUser(input FormUpdateUserInput) (User, error){
	user, err := s.repository.FindById(input.ID)

	if err != nil{
		return user, err
	}

	//mapping input ke tipe user
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	updatedUser, err := s.repository.Update(user)

	if err != nil{
		return updatedUser, err
	}

	return updatedUser, nil
}