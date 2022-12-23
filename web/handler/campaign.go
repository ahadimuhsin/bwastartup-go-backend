package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService	user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

// menampilkan data semua campaign
func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		//skip
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns" : campaigns})
}

func (h *campaignHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {

	}

	input := campaign.FormCreateCampaignInput{}
	input.Users = users
	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) Create(c *gin.Context){
	var input campaign.FormCreateCampaignInput

	err := c.ShouldBind(&input)

	if err != nil{
		users, e := h.userService.GetAllUsers()

		if e != nil{
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		input.Users = users
		input.Error = err

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}
	//ambil data user
	user, err := h.userService.GetUserById(input.UserID)
	//mapping form input
	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	
	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}