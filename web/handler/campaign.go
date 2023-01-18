package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService user.Service
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

//create page new campaign
func (h *campaignHandler) New(c *gin.Context){
	users, err := h.userService.GetAllUsers()

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := campaign.FormCreateCampaignInput{}
	input.Users = users
	c.HTML(http.StatusOK, "campaign_new.html", input)
}
//create page
func (h *campaignHandler) Create(c *gin.Context){
	var input campaign.FormCreateCampaignInput

	//binding
	err := c.ShouldBind(&input)

	//kalau error, tampilkan kembali halaman new campaign
	if err != nil {
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
	
	//get data user
	user, err := h.userService.GetUserById(input.UserID)

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	//mapping input
	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	//panggil service create campaign
	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) NewImage(c *gin.Context){
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam);

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": id})
}

func (h *campaignHandler) CreateImage(c *gin.Context){
	//tangkap id dan data campaignya
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam);
	existingCampaign, err := h.campaignService.GetCampaign(campaign.GetCampaignDetailInputById{ID: id})

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	//tangkap filenya
	file, err := c.FormFile("file")
	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	//set user id
	userID := existingCampaign.UserId
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignImageInput := campaign.CreateCampaignImageInput{}
	createCampaignImageInput.CampaignID = id
	createCampaignImageInput.IsPrimary = true
	

	//get data user
	userCampaign, err := h.userService.GetUserById(userID)
	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignImageInput.User = userCampaign

	_, err = h.campaignService.SaveCampaignImage(createCampaignImageInput, path)

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")

}
