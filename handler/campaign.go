package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"

	// "fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

// route api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	//tangkap parameter di handler
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userId)
	//kalau ada error
	if err != nil {
		errorMessage := gin.H{"errors": "Error occured when try to get list of campaigns"}
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//kalau berhasil
	response := helper.APIResponse("List of Campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// route api/v1/campaigns/v1
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	//get parameter route
	//repository : get campaign by slug

	var input campaign.GetCampaignDetailInput
	//untuk menangkap uri
	err := c.ShouldBindUri(&input)
	// fmt.Println(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//handler: mapping slug di url ke struct input -> service, call formatter
	//service: inputnya struct input -> menangkap slug di url
	campaignDetail, err := h.campaignService.GetCampaign(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign Detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	//tangkap parameter dari user ke input struct
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		//response menggunakna helper
		response := helper.APIResponse("Failed To Create Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//ambil current user dari jwt/handler
	//ambil data user dari Context gin, dari auth Middleware
	currentUser := c.MustGet("currentUser").(user.User)
	//set nilai idnya ke variabel
	input.User = currentUser
	//panggil service, parameternya input struct (dan buat slug)
	newCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil{
		response := helper.APIResponse("Failed to save campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	
	response := helper.APIResponse("Campaign created", http.StatusOK, "error", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}
