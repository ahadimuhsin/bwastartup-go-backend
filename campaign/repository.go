package campaign

import (
	// "fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Campaign, error) //pakai slice karena mengembalikan banyak data
	FindByUserId(id int) ([]Campaign, error)
	FindBySlug(Slug string) (Campaign, error)
	FindById(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign)(Campaign, error)
	CreateImage(campaignImage CampaignImage) (CampaignImage, error)
	CheckIsNonPrimary(campaignID int)(bool, error)
}

// bersifat private
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByUserId(id int) ([]Campaign, error) {
	var campaigns []Campaign //siapkan variable campaigns berbentuk slice
	//Preload untuk memanggil relasi, yg didefinisikan di model, dan memanggilnya dengan kondisi
	err := r.db.Where("user_id = ?", id).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (r *repository) FindBySlug(Slug string) (Campaign, error) {
	var campaign Campaign
	//Preload untuk memanggil relasi, yg didefinisikan di model, dan memanggilnya dengan kondisi
	err := r.db.Where("slug= ?", Slug).Preload("CampaignImages").Preload("User").Find(&campaign).Error
	// fmt.Println(Slug)
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

func (r *repository) FindById(id int) (Campaign, error) {
	var campaign Campaign
	//Preload untuk memanggil relasi, yg didefinisikan di model, dan memanggilnya dengan kondisi
	err := r.db.Where("id = ?", id).Preload("CampaignImages").Preload("User").Find(&campaign).Error
	// fmt.Println(Slug)
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

// menyimpan data campaign baru
func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

//update campaign
func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error

	if( err != nil){
		return campaign, err
	}

	return campaign, nil
}

func(r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error){
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}

func(r *repository) CheckIsNonPrimary(campaignID int)(bool, error){
	//query update
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}
	return true, nil
}