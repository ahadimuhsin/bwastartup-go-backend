package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
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
	//handler: mapping slug di url ke struct input -> service, call formatter
	//service: inputnya struct input -> menangkap slug di url
	//repository : get campaign by slug

	var input campaign.GetCampaignDetailInput
	//untuk menangkap uri
	err := c.ShouldBindUri(&input)
	// fmt.Println(input)
	if err != nil{
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.campaignService.GetCampaign(input)

	if err != nil{
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign Detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}
