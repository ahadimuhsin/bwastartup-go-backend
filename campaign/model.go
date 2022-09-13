package campaign

import (
	"bwastartup/user"
	"time"
)

//entity = model

type Campaign struct{
	ID 					int
	UserId				int
	Name 				string
	ShortDescription	string
	Description			string
	Perks				string
	BackerCount			int
	GoalAmount			int
	CurrentAmount		int
	Slug				string
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
	//untuk relasi
	CampaignImages		[]CampaignImage
	User 				user.User		
}

type CampaignImage struct{
	ID					int
	CampaignId			int
	FileName			string
	IsPrimary			int
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}