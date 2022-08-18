package user

import "gorm.io/gorm"

//perlu dibuat interface agar objek lain mengacu ke sini, tidak
//langsung ke implementasi
//bersifat public
type Repository interface {
	Save(user User) (User, error) //kurung pertama adalah nama variabel dan modelnya
	//kurung kedua adalah returnnya
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
