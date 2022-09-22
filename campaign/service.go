package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

// buat interface terlebih dahulu
type Service interface {
	//buat kontrak (list fungsinya nanti apa aja)
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaign(input GetCampaignDetailInputById) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputId GetCampaignDetailInputById, inputData CreateCampaignInput)(Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userId int) ([]Campaign, error) {
	//cek userid
	//jika user idnya tidak 0
	if userId != 0 {
		//tampilkan data dengan userId tersebut
		campaigns, err := s.repository.FindByUserId(userId)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	//jika tidak ada userId
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetCampaign(input GetCampaignDetailInputById) (Campaign, error) {
	campaign, err := s.repository.FindById(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UserId = input.User.ID

	campaignSlug := fmt.Sprintf("%s-%d", input.Name, input.User.ID)
	//proses pembuatan slug
	campaign.Slug = slug.Make(campaignSlug)

	//proses penyimpanan ke db
	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil

}

func (s *service) UpdateCampaign(inputId GetCampaignDetailInputById, inputData CreateCampaignInput)(Campaign, error){
	campaign, err := s.repository.FindById(inputId.ID)

	if err != nil{
		return campaign, err
	}

	if campaign.UserId != inputData.User.ID{
		return campaign, errors.New("Not an owner of the campaign")
	}
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount
	// campaignSlug := fmt.Sprintf("%s-%d", inputData.Name, inputData.User.ID)
	// //proses pembuatan slug
	// campaign.Slug = slug.Make(campaignSlug)
	//lakukan pengecekan campaign ini punya user siapa
	
	//proses update
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil{
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}