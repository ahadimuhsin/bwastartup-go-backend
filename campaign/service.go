package campaign

//buat interface terlebih dahulu
type Service interface{
	//buat kontrak (list fungsinya nanti apa aja)
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaign(input GetCampaignDetailInput)(Campaign, error)
}

type service struct{
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userId int) ([]Campaign, error){
	//cek userid
	//jika user idnya tidak 0
	if userId != 0{
		//tampilkan data dengan userId tersebut
		campaigns, err := s.repository.FindByUserId(userId)
		if err != nil{
			return campaigns, err
		}
		return campaigns, nil
	}

	//jika tidak ada userId
	campaigns, err := s.repository.FindAll()
	if err != nil{
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetCampaign(input GetCampaignDetailInput) (Campaign, error){
	campaign, err := s.repository.FindBySlug(input.Slug)
	if err != nil{
		return campaign, err
	}
	return campaign, nil
}