package user

import "gorm.io/gorm"

//perlu dibuat interface agar objek lain mengacu ke sini, tidak
//langsung ke implementasi
//bersifat public
type Repository interface {
	Save(user User) (User, error) //kurung pertama adalah nama variabel dan modelnya
	//kurung kedua adalah returnnya
	FindByEmail(email string) (User, error)
	FindById(id int) (User, error)
	Update(user User) (User, error)
	FindAll() ([]User, error)
}

//bersifat private
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	//membuat user baru
	err := r.db.Create(&user).Error

	//cek error
	if err!=nil{
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	//cari user dengan email tertentu
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if( err != nil){
		return user, err
	}

	return user, nil
}

func (r *repository) FindById(id int) (User, error) {
	//cari user dengan id tertentu
	var user User
	err := r.db.Where("id = ?", id).Find(&user).Error

	if( err != nil){
		return user, err
	}

	return user, nil
}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if( err != nil){
		return user, err
	}

	return user, nil
}

func (r *repository) FindAll() ([]User, error){
	var users []User

	err := r.db.Find(&users).Error
	if err != nil{
		return users, err
	}

	return users, nil
}