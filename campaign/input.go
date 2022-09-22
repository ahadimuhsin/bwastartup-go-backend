package campaign

import "bwastartup/user"

//untuk mengambil input dari user
type GetCampaignDetailInput struct{
	Slug string `uri:"slug" binding:"required"`
}

type GetCampaignDetailInputById struct{
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct{
	Name 				string `json:"name" binding:"required"`
	ShortDescription 	string `json:"short_description" binding:"required"`
	GoalAmount 			int `json:"goal_amount" binding:"required"`
	Description 		string `json:"description" binding:"required"`
	Perks 				string `json:"perks" binding:"required"`
	User				user.User
}