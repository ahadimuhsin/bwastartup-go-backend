package campaign

import (
	"bwastartup/user"
	"time"

	"github.com/leekchan/accounting"
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
//memformat goal amount dari struct campaign
func (c Campaign) GoalAmountFormat() string{
	ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}

	return ac.FormatMoney(c.GoalAmount)
}

type CampaignImage struct{
	ID					int
	CampaignId			int
	FileName			string
	IsPrimary			int
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}