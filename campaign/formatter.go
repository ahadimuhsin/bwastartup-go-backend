package campaign

type CampaignFormatter struct{
	ID 					int `json:"id"`
	UserId 				int `json:"user_id"`
	Name 				string `json:"name"`
	ShortDescription 	string `json:"short_description"`
	ImageUrl 			string `json:"image_url"`
	GoalAmount 			int `json:"goal_amount"`
	CurrentAmount 		int `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter{
	var ImageUrl string
	if len(campaign.CampaignImages) > 0{
		ImageUrl = campaign.CampaignImages[0].FileName
	}
	campaignFormatter := CampaignFormatter{
		ID : campaign.ID,
		UserId : campaign.UserId,
		ShortDescription : campaign.ShortDescription,
		GoalAmount : campaign.GoalAmount,
		CurrentAmount : campaign.CurrentAmount,
		ImageUrl: ImageUrl,
	}
	return campaignFormatter
}

//handle slice of campaigns to use formatter
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter{
	if len(campaigns) == 0{
		return []CampaignFormatter{}
	}
	var campaignsFormatter []CampaignFormatter
	for _, campaign := range campaigns{
		singleCampaignFormatter:= FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, singleCampaignFormatter)
	}
	return campaignsFormatter
}