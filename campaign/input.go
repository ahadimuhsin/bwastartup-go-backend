package campaign

//untuk mengambil input dari user
type GetCampaignDetailInput struct{
	Slug string `uri:"slug" binding:"required"`
}